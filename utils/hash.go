package utils

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// HashPassword hashes a plain password using bcrypt.
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return "", err
	}
	return string(hashedBytes), nil
}

// ComparePasswords checks if the plain password matches the hashed password.
func ComparePasswords(plainPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
