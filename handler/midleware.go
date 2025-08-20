package handler

import (
	"database/sql"
	"errors"
	"github/http/copy/task4/config"
	c "github/http/copy/task4/constants"
	"github/http/copy/task4/security"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/cast"
)

func (h *Handler) AuthMiddleware(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(c.UnAuth, gin.H{"error": "unauthorized"})
		ctx.Abort()
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		ctx.JSON(c.BadReq, gin.H{"error": "invalid token format"})
		ctx.Abort()
		return
	}

	var claims security.Claims
	token, err := jwt.ParseWithClaims(parts[1], &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return config.Cfg().JWTsecretkey, nil
	})

	if err != nil || !token.Valid {
		ctx.JSON(c.UnAuth, gin.H{"error": "unauthorized"})
		ctx.Abort()
		return
	}

	data, err := h.DB.Auth().GetSessionId(ctx)
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
	claims.ID = data.ID

	session, err := h.DB.Auth().GetSessionByID(ctx, cast.ToInt(claims.ID))
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(c.UnAuth, gin.H{"error": "unauthorized"})
		ctx.Abort()
		return
	} else if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	if !session.ExpiredAt.After(time.Now()) {
		ctx.JSON(c.UnAuth, gin.H{"error": "token expired"})
		ctx.Abort()
		return
	}

	ctx.Set("userId", claims.UserId)
	ctx.Set("role", claims.Role)

	fullPath := ctx.FullPath()
	if fullPath == "" {
		fullPath = ctx.Request.URL.Path
	}
	if strings.HasPrefix(fullPath, "/admin") && claims.Role != "admin" {
		ctx.JSON(c.Forbidden, gin.H{"error": "forbidden"})
		ctx.Abort()
		return
	}

	ctx.Next()
}
