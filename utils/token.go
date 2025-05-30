package utils

import (
	"errors"
	"time"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Claims struct for custom payload
type Claims struct {
	UserId string `json:"userId"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token with a 1-hour expiry
func GenerateToken(userId, role string) (string, error) {
	claims := &Claims{
		UserId: userId,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(6 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// VerifyToken verifies and parses a JWT token
func VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("could not parse claims")
	}
	return claims, nil
}

func VerifyTokenWithRoles(tokenString string, allowedRoles []string) (*Claims, error) {
	claims, err := VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	for _, role := range allowedRoles {
		if claims.Role == role {
			return claims, nil
		}
	}
	return nil, errors.New("user does not have required privileges")
}