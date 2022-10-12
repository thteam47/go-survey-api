package models

import (
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	TokenInfo *TokenInfo `json:"token_info,omitempty" bson:"token_info,omitempty"`
	jwt.StandardClaims
}
