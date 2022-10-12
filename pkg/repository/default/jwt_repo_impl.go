package repoimpl

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/thteam47/go-identity-authen-api/errutil"
	"github.com/thteam47/go-identity-authen-api/pkg/db"
	"github.com/thteam47/go-identity-authen-api/pkg/models"
	"github.com/thteam47/go-identity-authen-api/pkg/repository"
)

type JwtRepositoryImpl struct {
	handler *db.Handler
}


func NewJwtRepo(handler *db.Handler) repository.JwtRepository {
	return &JwtRepositoryImpl{
		handler: handler,
	}
}
func (inst *JwtRepositoryImpl) Generate(tokenInfo *models.TokenInfo) (string, error) {
	claims := models.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(tokenInfo.Exp),
		},
		TokenInfo: tokenInfo,
	}

	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := claim.SignedString([]byte(inst.handler.JwtKey))
	if err != nil {
		return "", errutil.Wrapf(err, "claim.SignedString")
	}
	return token, nil
}
func (inst *JwtRepositoryImpl) Verify(accessToken string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&models.Claims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(inst.handler.JwtKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil

}
