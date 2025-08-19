package security

import (
	"github/http/copy/task4/config"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserId int    `json:"userId"`
	Role   string `json:"role"`
	ID     int    `json:"id"`
	jwt.StandardClaims
}

func GenerateJWTToken(userId, id int, role string) (string, error) {
	expiredTime := jwt.TimeFunc().Add(time.Hour * 24)

	if expiredTime.Before(time.Now()) || expiredTime.Equal(time.Now()) {

	}
	claims := Claims{
		UserId: userId,
		Role:   role,
		ID:     id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(config.Cfg().JWTsecretkey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
