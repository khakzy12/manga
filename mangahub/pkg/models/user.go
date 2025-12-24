package models

import "github.com/golang-jwt/jwt/v5"

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
