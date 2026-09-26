package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"todo_api/internal/config"
	"todo_api/internal/store"

	"github.com/gin-gonic/gin"
)

func bearerFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if after, ok := strings.CutPrefix(authHeader, "Bearer "); ok {
		return after
	}
	return ""
}

func AuthMiddleware(cfg *config.Config, r *store.Redis) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("access_token")
		if tokenStr == "" {
			tokenStr = bearerFromHeader(c)
		}

		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing header"})
			return
		}

		claims, err := ParseAccess(tokenStr, cfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		ctx := context.Background()
		if _, err := r.GetUserByJTI(ctx, "access: "+claims.ID); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
			return
		}

		c.Set("user_id", claims.Subject)
		c.Next()

	}
}

func MustCookie(c *gin.Context, name string) (string, error) {
	val, err := c.Cookie(name)
	if err != nil || val == "" {
		return "", errors.New("missing cookie: " + name)
	}

	return val, nil
}
