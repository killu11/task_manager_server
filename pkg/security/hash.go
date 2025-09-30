package security

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	HashPassErr   = errors.New("hash_failed")
	VerifyPassErr = errors.New("invalid_password")
)

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashPassword []byte, inputPass string) bool {
	err := bcrypt.CompareHashAndPassword(hashPassword, []byte(inputPass))
	return err == nil
}
