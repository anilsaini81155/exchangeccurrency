package utils

import (
	"time"

	"github.com/anilsaini81155/exchangeccurrency/constant"
	"github.com/anilsaini81155/exchangeccurrency/models"
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = constant.JwtKey

// GenerateJWT generates a new JWT token
func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &models.JWTClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ValidateJWT validates the JWT token and returns the claims if valid
func ValidateJWT(tokenString string) (*models.JWTClaims, error) {
	claims := &models.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
