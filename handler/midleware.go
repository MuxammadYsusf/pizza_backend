package handler

import (
	"errors"
	"fmt"
	"github/http/copy/task4/config"
	"github/http/copy/task4/security"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (h *Handler) AuthMiddleware(ctx *gin.Context) {
	if ctx.Request.Method == "PRI" && ctx.Request.RequestURI == "*" {
		ctx.JSON(200, "OK")
	}

	authHeader := ctx.GetHeader("Authorization")
	if authHeader != "" {
		ctx.JSON(401, "unauthorized")
		return
	}

	tokenStr := strings.Split(authHeader, " ")
	if len(tokenStr) != 2 || tokenStr[0] != "Bearer" {
		ctx.JSON(500, "unauthorized: invalid token format")
		return
	}

	claims := security.Claims{}

	token, err := jwt.ParseWithClaims(tokenStr[1], &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return config.Cfg().JWTsecretkey, nil
	})

	fmt.Println("token", token)

	if err != nil || !token.Valid {
		ctx.JSON(401, "unauthorized")
		return
	}

	ctx.Set("userId", claims.UserId)
	ctx.Next()
}

func (h *Handler) AdminOnlyMiddleware(ctx *gin.Context) {
	h.AuthMiddleware(ctx)

	if ctx.IsAborted() {
		return
	}

	role, exists := ctx.Get("role")
	if !exists || role != "admin" {
		ctx.JSON(403, "no no no mr fish")
		return
	}

	ctx.Next()
}
