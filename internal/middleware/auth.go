package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"
	"todo_api/internal/config"
	"todo_api/internal/models"
	"todo_api/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func IssueTokens(userID string, cfg *config.Config) (*models.Tokens, error) {
	now := time.Now().UTC()
	t := &models.Tokens{
		UserID:   userID,
		JTIAcc:   uuid.NewString(),
		JTIRef:   uuid.NewString(),
		ExpAcc:   now.Add(15 * time.Minute),
		ExpRef:   now.Add(7 * 24 * time.Hour),
		Issuer:   "todo-jwt-app",
		Audience: "todo-jwt-client",
	}

	acc := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ID:        t.JTIAcc,
		Issuer:    t.Issuer,
		Audience:  jwt.ClaimStrings{t.Audience},
		ExpiresAt: jwt.NewNumericDate(t.ExpAcc),
		IssuedAt:  jwt.NewNumericDate(now),
	})

	ref := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ID:        t.JTIRef,
		Issuer:    t.Issuer,
		Audience:  jwt.ClaimStrings{t.Audience},
		ExpiresAt: jwt.NewNumericDate(t.ExpRef),
		IssuedAt:  jwt.NewNumericDate(now),
	})

	var err error
	t.Access, err = acc.SignedString([]byte(cfg.AccessSecret))
	if err != nil {
		return nil, err
	}

	t.Refresh, err = ref.SignedString([]byte(cfg.RefreshSecret))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func Persist(ctx context.Context, r *store.Redis, t *models.Tokens) error {
	if err := r.SetJTI(ctx, "access: "+t.JTIAcc, t.UserID, t.ExpAcc); err != nil {
		return err
	}

	if err := r.SetJTI(ctx, "refresh: "+t.JTIRef, t.UserID, t.ExpRef); err != nil {
		return err
	}

	return nil
}

func SetAuthCookies(c *gin.Context, t *models.Tokens) {
	c.SetSameSite(http.SameSiteLaxMode)
	//domain should be the website name e.g "td-app.com" and secure set to true in production
	c.SetCookie("access_token", t.Access, int(time.Until(t.ExpAcc).Seconds()), "/", "localhost", false, true)
	c.SetCookie("refresh_token", t.Refresh, int(time.Until(t.ExpRef).Seconds()), "/", "localhost", false, true)
}

func ClearAuthCookies(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	//domain should be the website name e.g "td-app.com" and secure set to true in production
	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("refresh-token", "", -1, "/", "localhost", false, true)
}

func ParseAccess(tokenStr string, cfg *config.Config) (*jwt.RegisteredClaims, error) {
	return parseWithSecret(tokenStr, cfg.AccessSecret)
}

func ParseRefresh(tokenStr string, cfg *config.Config) (*jwt.RegisteredClaims, error) {
	return parseWithSecret(tokenStr, cfg.RefreshSecret)
}

func parseWithSecret(tokenStr string, secret string) (*jwt.RegisteredClaims, error) {
	if secret == "" {
		return nil, errors.New("jwt secret not configured")
	}

	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithAudience("todo-jwt-client"),
		jwt.WithIssuer("todo-jwt-app"),
	)

	token, err := parser.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
