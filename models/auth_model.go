package models

import (
	"github.com/golang-jwt/jwt/v5"
)

// Auth model for sign in
type Auth struct {
	Username string `json:"username,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

// Claims model for token
type Claims struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
