package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func GetHashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassowrdValid(rawPassword, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))
	if err == nil {
		return true, nil
	} else if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	} else {
		return false, err
	}
}
