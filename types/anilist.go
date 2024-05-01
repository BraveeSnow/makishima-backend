package types

type AnilistTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type AnilistIdentity struct {
	Viewer struct {
		ID int `json:"id"`
	}
}
