package handler

import (
	"errors"
	"net/http"
	"strings"

	"github/http/copy/task4/pkg/security"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// AuthMiddleware проверяет Bearer JWT, кладёт userId и role в контекст.
func (h *Handler) AuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		abortUnauthorized(c, "missing Authorization header")
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		abortUnauthorized(c, "invalid Authorization header format")
		return
	}

	var claims security.Claims
	token, err := jwt.ParseWithClaims(parts[1], &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(h.cfg.JWTSecret), nil
	})
	if err != nil || token == nil || !token.Valid {
		abortUnauthorized(c, "invalid token")
		return
	}

	c.Set("userId", claims.UserId)
	if claims.Role != "" {
		c.Set("role", claims.Role)
	}
	c.Next()
}

// AdminOnlyMiddleware должен подключаться ПОСЛЕ AuthMiddleware в цепочке.
// Пример: group := r.Group("/admin", h.AuthMiddleware, h.AdminOnlyMiddleware)
func (h *Handler) AdminOnlyMiddleware(c *gin.Context) {
	roleVal, ok := c.Get("role")
	if !ok {
		abortForbidden(c, "role missing")
		return
	}
	role, _ := roleVal.(string)
	if role != "admin" {
		abortForbidden(c, "admin only")
		return
	}
	c.Next()
}

/* ---------- helpers ---------- */

func abortUnauthorized(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error":       "unauthorized",
		"description": msg,
	})
}

func abortForbidden(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"error":       "forbidden",
		"description": msg,
	})
}
