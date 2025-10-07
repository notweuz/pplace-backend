package model

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}
