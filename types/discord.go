package types

import "github.com/golang-jwt/jwt/v5"

type DiscordIdentity struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type DiscordJwtClaims struct {
	jwt.RegisteredClaims
	DiscordId string `json:"id"`
}

type DiscordTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}
