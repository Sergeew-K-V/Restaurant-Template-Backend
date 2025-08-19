package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pass string) (string, error) {
	hash, error := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	if error != nil {
		fmt.Println("Failed to hash password", error)
	}

	return string(hash), nil
}

func CheckPassword(hash, pass string) bool {
	error := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))

	return error == nil
}
