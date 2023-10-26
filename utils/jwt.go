package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Data string `json:"data"`
	jwt.RegisteredClaims
}

type Tokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func getSecretKey() []byte {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		secretKey = "random-secret-123"
	}

	return []byte(secretKey)
}

func getTokenExpiration() (time.Duration, time.Duration) {
	accessTokenExp := os.Getenv("ACCESS_TOKEN_EXP")
	refreshTokenExp := os.Getenv("REFRESH_TOKEN_EXP")

	accessTokenExpiration, err := strconv.Atoi(accessTokenExp)
	if err != nil {
		accessTokenExpiration = 7
	}

	refreshTokenExpiration, err := strconv.Atoi(refreshTokenExp)
	if err != nil {
		refreshTokenExpiration = 14
	}

	return time.Hour * 24 * time.Duration(accessTokenExpiration), time.Hour * 24 * time.Duration(refreshTokenExpiration)
}

func GenerateTokens(userID string) (Tokens, error) {
	accessTokenExp, refreshTokenExp := getTokenExpiration()
	secretkey := getSecretKey()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Data: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExp)),
		},
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Data: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExp)),
		},
	})

	accessTokenString, err := accessToken.SignedString(secretkey)

	if err != nil {
		return Tokens{}, err
	}

	refreshTokenString, err := refreshToken.SignedString(secretkey)

	if err != nil {
		return Tokens{}, err
	}

	return Tokens{
		Access:  accessTokenString,
		Refresh: refreshTokenString,
	}, nil
}

func VerifyToken(tokenString string) (*Claims, error) {
	secretKey := getSecretKey()

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token string")
}

func ReIssueAccessToken(refreshToken string) (Tokens, error) {
	claims, err := VerifyToken(refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	if time.Unix(claims.ExpiresAt.Unix(), 0).After(time.Now()) {
		return GenerateTokens(claims.Data)
	}

	return Tokens{}, fmt.Errorf("refresh token has expired")
}
