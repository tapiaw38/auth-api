package models

import "github.com/golang-jwt/jwt"

// AppClaims is the model for the claims
type AppClaims struct {
	UserId string `json:"userId"`
	Email  string `json:"email"`
	jwt.StandardClaims
}
