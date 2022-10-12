package util

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func GenerateTotp(secret string) (*otp.Key, error) {
	return totp.Generate(totp.GenerateOpts{
		Issuer:      "ThteamM",
		AccountName: "thteam47@gmail.com",
		Secret:      []byte(secret),
	})
}
