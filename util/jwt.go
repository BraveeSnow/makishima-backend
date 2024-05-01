package util

import (
	"fmt"
	"makishima-backend/types"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(identity *types.DiscordIdentity, expiry int64) (string, error) {
	currentTime := time.Unix(time.Now().Unix(), 0)
	claims := types.DiscordJwtClaims{
		DiscordId: identity.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(expiry) * time.Second)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("MAKISHIMA_SIGKEY")))
}

func DecodeToken(token *string) (*types.DiscordJwtClaims, error) {
	claims := types.DiscordJwtClaims{}
	_, err := jwt.ParseWithClaims(*token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("MAKISHIMA_SIGKEY")), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token is invalid - %s", err.Error())
	}

	return &claims, nil
}
