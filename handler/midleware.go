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
		ctx.JSON(c.UnAuth, "unauthorized")
		return
	}

	tokenStr := strings.Split(authHeader, " ")
	if len(tokenStr) != 2 || tokenStr[0] != "Bearer" {
		ctx.JSON(c.BadReq, gin.H{"unauthorized": " invalid token format"})
		return
	}

	expiredTime := jwt.TimeFunc().Add(time.Hour * 24)

	claims := security.Claims{}

	if expiredTime.Before(time.Now()) || expiredTime.Equal(time.Now()) {
		ctx.JSON(c.UnAuth, gin.H{"error": "unauthorized"})
		return
	}

	token, err := jwt.ParseWithClaims(tokenStr[1], &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return config.Cfg().JWTsecretkey, nil
	})

	if err != nil || !token.Valid {

		ctx.JSON(c.UnAuth, gin.H{"error": "unauthorized"})
		return
	}

	data, err := h.DB.Auth().GetSessionId(ctx)
	if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return
	}

	claims.ID = data.ID

	session, err := h.DB.Auth().GetSessionByID(ctx, cast.ToInt(claims.ID))
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(c.UnAuth, gin.H{"error": "unauthorized"})
		return
	} else if err != nil {
		ctx.JSON(c.Err, gin.H{"error": err.Error()})
		return

	}

	if session.ExpiredAt.Before(time.Now()) || session.ExpiredAt.Equal(time.Now()) {
		ctx.JSON(c.UnAuth, gin.H{"error": "token expired"})
		return
	}

	role := claims.Role

	if role != "admin" {
		ctx.JSON(c.Forbidden, gin.H{"error": errors.New("no no no mr fish")})
		return
	}

	ctx.Set("userId", claims.UserId)
	ctx.Next()
}
