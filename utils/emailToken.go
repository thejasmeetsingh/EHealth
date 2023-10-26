package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GetToken(duration time.Time, data string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Data: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(duration),
		},
	})

	signedToken, err := token.SignedString([]byte(getSecretKey()))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
