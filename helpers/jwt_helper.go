package helpers

import "github.com/golang-jwt/jwt/v4"

type JwtCustomClaims struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}
