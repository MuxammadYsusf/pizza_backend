package security

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserId int    `json:"userId"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func GenerateJWT(m map[string]interface{}, tokenExpireTime time.Duration, tokenSecretKey string) (tokenString string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	for key, value := range m {
		claims[key] = value
	}

	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(tokenExpireTime).Unix()

	tokenString, err = token.SignedString([]byte(tokenSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractToken(bearer string) (token string, err error) {
	strArr := strings.Split(bearer, " ")
	if len(strArr) == 2 {
		return strArr[1], nil
	}
	return token, errors.New("wrong token format")
}
