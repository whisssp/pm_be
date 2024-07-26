package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordBytes), nil
}

func ComparePasswords(hashed []byte, plain []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashed, plain)
	return err == nil
}