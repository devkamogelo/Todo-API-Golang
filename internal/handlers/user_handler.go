package handlers

import (
	"context"
	"net/http"
	"strings"
	"todo_api/internal/config"
	"todo_api/internal/middleware"
	"todo_api/internal/models"
	"todo_api/internal/repository"
	"todo_api/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var registerRequest RegisterRequest

		if err := c.BindJSON(&registerRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(registerRequest.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters long"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password " + err.Error()})
			return
		}

		user := &models.User{
			Email:    registerRequest.Email,
			Password: string(hashedPassword),
		}

		createdUser, err := repository.CreateUser(pool, user)

		if err != nil {
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, createdUser)
	}
}

func LoginHandler(pool *pgxpool.Pool, cfg *config.Config, rds *store.Redis) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest LoginRequest

		err := c.BindJSON(&loginRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := repository.GetUserByEmail(pool, loginRequest.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Credentials"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Credentials"})
			return
		}

		tokens, err := middleware.IssueTokens(user.ID, cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not issue token"})
			return
		}

		if err = middleware.Persist(c, rds, tokens); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not persist token"})
			return
		}

		middleware.SetAuthCookies(c, tokens)

		c.JSON(http.StatusOK, true)
	}
}

func RefreshHandler(cfg *config.Config, rds *store.Redis) gin.HandlerFunc {
	return func(c *gin.Context) {
		ref, err := middleware.MustCookie(c, "refresh_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims, err := middleware.ParseRefresh(ref, cfg)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx := context.Background()
		if _, err := rds.GetUserByJTI(ctx, "refresh: "+claims.ID); err != nil {
			c.JSON(http.StatusUnauthorized, "token revoked")
			return
		}

		_ = rds.DelJTI(ctx, "refresh: "+claims.ID)

		tokens, err := middleware.IssueTokens(claims.Subject, cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not issue token"})
			return
		}

		if err = middleware.Persist(c, rds, tokens); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not persist token"})
			return
		}

		middleware.SetAuthCookies(c, tokens)

		c.JSON(http.StatusOK, true)

	}
}

func LogoutHandler(cfg *config.Config, rds *store.Redis) gin.HandlerFunc {
	return func(c *gin.Context) {
		acc, _ := c.Cookie("access_token")
		ref, _ := c.Cookie("refresh_token")
		ctx := context.Background()

		if acc != "" {
			if claims, err := middleware.ParseAccess(acc, cfg); err != nil {
				_ = rds.DelJTI(ctx, claims.ID)
			}
		}

		if ref != "" {
			if claims, err := middleware.ParseRefresh(ref, cfg); err != nil {
				_ = rds.DelJTI(ctx, claims.ID)
			}
		}

		middleware.ClearAuthCookies(c)
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}

}

func TestProtectedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "protected route accessed successfully",
			"user_id": userID,
		})
	}
}
