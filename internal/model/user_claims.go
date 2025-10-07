package model

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	ID           uint `json:"id"`
	TokenVersion uint `json:"token_version"`
	jwt.RegisteredClaims
}
