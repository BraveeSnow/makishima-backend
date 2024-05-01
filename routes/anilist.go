package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"makishima-backend/types"
	"makishima-backend/types/database"
	"makishima-backend/util"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

const AnilistAPIEndpoint = "https://anilist.co/api/v2"
const AnilistGQLEndpoint = "https://graphql.anilist.co"

const AnilistGQLIdentify = `
	query {
		Viewer {
			id
		}
	}
`

func AnilistOAuthRedirect(data *types.MakishimaData) gin.HandlerFunc {
	anilistOAuthUrl, _ := url.JoinPath(AnilistAPIEndpoint, "/oauth/token")

	makishimaPanelUrl, err := url.JoinPath(os.Getenv("MAKISHIMA_URI"), "/panel")
	if err != nil {
		data.Logger.Fatal().Msgf("invalid URI set in 'MAKISHIMA_URI'")
	}

	anilistId := os.Getenv("ANILIST_ID")
	anilistSecret := os.Getenv("ANILIST_SECRET")
	anilistRedirect := os.Getenv("ANILIST_REDIRECT")

	identifyQuery, _ := json.Marshal(map[string]string{
		"query": AnilistGQLIdentify,
	})

	return func(ctx *gin.Context) {
		identity, err := ctx.Cookie("identity")

		if err != nil {
			data.Logger.Warn().Msg("no identity was given to associate with")
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "no identity provided",
			})
			return
		}

		claims, err := util.DecodeToken(&identity)
		if err != nil {
			data.Logger.Warn().Msg(err.Error())
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid identity",
			})
			return
		}

		// find requesting user
		user := database.User{}
		result := data.Database.Find(&user, claims.DiscordId)
		if result.Error != nil {
			data.Logger.Error().Msgf("valid JWT was passed but encountered a database error - %s", result.Error.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "unknown",
			})
			return
		}

		// process oauth request
		code, codeExists := ctx.GetQuery("code")
		if !codeExists {
			data.Logger.Warn().Msg("no code given")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request",
			})
			return
		}

		formBody := map[string]string{
			"grant_type":    "authorization_code",
			"client_id":     anilistId,
			"client_secret": anilistSecret,
			"redirect_uri":  anilistRedirect,
			"code":          url.QueryEscape(code),
		}
		formJson, _ := json.Marshal(formBody)
		println(string(formJson))
		tokenRequest, _ := http.NewRequest("POST", anilistOAuthUrl, bytes.NewReader(formJson))
		tokenRequest.Header.Add("Content-Type", "application/json")
		tokenRequest.Header.Add("Accept", "application/json")
		tokenResponse, err := http.DefaultClient.Do(tokenRequest)

		if err != nil {
			data.Logger.Error().Msgf("unable to fetch anilist tokens - %s", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "unknown",
			})
			return
		}

		tokenBytes, err := io.ReadAll(tokenResponse.Body)
		if err != nil {
			data.Logger.Error().Msgf("unable to read response - %s", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "unknown",
			})
			return
		}

		if tokenResponse.StatusCode != 200 {
			data.Logger.Error().Msgf("token response returned code status code %d:\n%s", tokenResponse.StatusCode, string(tokenBytes))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "unknown",
			})
			return
		}

		tokensDeserialized := types.AnilistTokenResponse{}
		if json.Unmarshal(tokenBytes, &tokensDeserialized) != nil {
			data.Logger.Error().Msg("received malformed JSON response")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "unknown",
			})
			return
		}
		println(tokensDeserialized.AccessToken)

		// fetch anilist user ID
		identityRequest, _ := http.NewRequest("POST", AnilistGQLEndpoint, bytes.NewReader(identifyQuery))
		identityRequest.Header.Add("Content-Type", "application/json")
		identityRequest.Header.Add("Accept", "application/json")
		identityRequest.Header.Add("Authorization", "Bearer "+tokensDeserialized.AccessToken)
		identityResponse, err := http.DefaultClient.Do(identityRequest)

		if err != nil {
			data.Logger.Error().Msgf("AniList OAuth token exchange failed - %s", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "unknown",
			})
			return
		}

		identityBody, err := io.ReadAll(identityResponse.Body)
		if err != nil {
			data.Logger.Error().Msgf("unable to read identity response")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "unknown",
			})
			return
		}

		if identityResponse.StatusCode != 200 {
			data.Logger.Error().Msgf("identity response returned status code %d:\n%s", identityResponse.StatusCode, string(identityBody))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "unknown",
			})
			return
		}

		identityDeserialized := types.AnilistIdentity{}
		if json.Unmarshal(identityBody, &identityDeserialized) != nil {
			data.Logger.Error().Msg("unable to deserialize AniList identity - malformed JSON response")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "unknown",
			})
			return
		}

		// store tokens into database
		anilistEntry := database.AnilistUser{
			ID:           identityDeserialized.Viewer.ID,
			AccessToken:  tokensDeserialized.AccessToken,
			RefreshToken: tokensDeserialized.RefreshToken,
			TokenExpiry:  time.Now().Unix() + tokensDeserialized.ExpiresIn,
			UserID:       claims.DiscordId,
		}
		dbResult := data.Database.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"access_token", "refresh_token", "token_expiry"}),
		}).Create(&anilistEntry)

		if dbResult.Error != nil {
			data.Logger.Error().Msgf("unable to insert into database - %s", dbResult.Error.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "unknown",
			})
			return
		}

		// redirect to panel
		ctx.Redirect(http.StatusTemporaryRedirect, makishimaPanelUrl)
	}
}
