package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the JWT claims structure
type JWTClaims struct {
	UserID uint   `json:"id"` // Strapi uses "id" instead of "user_id"
	Email  string `json:"email,omitempty"`
	jwt.RegisteredClaims
}

// ExtractBearerToken extracts the bearer token from Authorization header
func ExtractBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	return parts[1], nil
}

// ParseJWTToken parses and validates a JWT token
func ParseJWTToken(tokenString string) (*JWTClaims, error) {
	// Get JWT secret from environment
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key" // fallback for development
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ExtractUserIDFromToken extracts user ID from Authorization header
func ExtractUserIDFromToken(authHeader string) (uint, error) {
	tokenString, err := ExtractBearerToken(authHeader)
	if err != nil {
		return 0, err
	}

	claims, err := ParseJWTToken(tokenString)
	if err != nil {
		return 0, err
	}

	return claims.UserID, nil
}
