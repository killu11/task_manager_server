package security

import (
	"errors"
	"fmt"
	"os"
	"time"
	_ "time"

	"github.com/golang-jwt/jwt"
	"github.com/goloop/env"
)

var (
	GenerateJWTErr = errors.New("jwt_generate_failed")
	SecretKeyErr   = errors.New("empty_secret_key")
)

type JWTClaims struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	jwt.StandardClaims
}

func newJWTClaims(id int, name string) JWTClaims {
	return JWTClaims{
		ID:       id,
		Username: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // Время жизни
			IssuedAt:  time.Now().Unix(),                     // Время создания
			Issuer:    "task_manager",
			NotBefore: time.Now().Unix(), // С какого момента действителен
		},
	}
}

func GenerateJWT(id int, name string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		newJWTClaims(id, name),
	)
	key := env.Get("SECRET_KEY")

	if key == "" {
		return "", fmt.Errorf("secret key missing: %v", SecretKeyErr)
	}

	signedToken, err := token.SignedString([]byte(key))

	if err != nil {
		return "", fmt.Errorf("jwt signature err: %v", err)
	}

	return signedToken, nil
}

func ParseUserID(tokenString string) (int, error) {
	claims := new(JWTClaims)
	_, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			secretKey := os.Getenv("SECRET_KEY")

			if secretKey == "" {
				return nil, SecretKeyErr
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

	if err != nil {
		return 0, err
	}
	return claims.ID, nil
}

func CheckTokenValid(tokenString string) (bool, error) {
	key := env.Get("SECRET_KEY")
	if key == "" {
		return false, fmt.Errorf("secret key missing: %v", SecretKeyErr)
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		new(JWTClaims),
		func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})

	if err != nil {
		return false, fmt.Errorf("failed parse token: %v", err)
	}
	return token.Valid, nil
}
