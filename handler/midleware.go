package handler

import (
	"errors"
	"fmt"
	"github/http/copy/task4/config"
	c "github/http/copy/task4/constants"
	"github/http/copy/task4/security"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (h *Handler) AuthMiddleware(ctx *gin.Context) {

	fmt.Println("1")
	if ctx.Request.Method == "PRI" && ctx.Request.RequestURI == "*" {
		ctx.JSON(c.OK, "")
	}

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

	claims := security.Claims{}

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

	ctx.Set("userId", claims.UserId)
	ctx.Next()
}

func (h *Handler) AdminOnlyMiddleware(ctx *gin.Context) {
	h.AuthMiddleware(ctx)

	fmt.Println("11")

	if ctx.IsAborted() {
		return
	}

	tokenStr := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")
	token, _, _ := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	role := token.Claims.(jwt.MapClaims)["role"].(string)

	if role != "admin" {
		ctx.JSON(c.Forbidden, gin.H{"error": errors.New("no no no mr fish")})
		return
	}

	ctx.Next()
}
