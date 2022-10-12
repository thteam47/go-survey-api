package repository

import "github.com/thteam47/go-identity-authen-api/pkg/models"

type JwtRepository interface {
	Generate(tokenInfo *models.TokenInfo) (string, error)
	Verify(accessToken string) (*models.Claims, error)
}
