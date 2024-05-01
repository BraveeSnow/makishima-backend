package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"makishima-backend/types"
	"makishima-backend/types/database"
	"makishima-backend/util"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

const DiscordEndpoint = "https://discord.com/api/v10"

func identifyDiscordUser(token *string) (*types.DiscordIdentity, error) {
	identifyUrl, _ := url.JoinPath(DiscordEndpoint, "/users/@me")
	request, _ := http.NewRequest("GET", identifyUrl, nil)
	request.Header.Add("Authorization", "Bearer "+*token)
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, fmt.Errorf("unable to perform identify request - %s", err.Error())
	}

	body, err := io.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("unable to identify user - status code %d:\n%s", response.StatusCode, string(body))
	}

	if err != nil {
		return nil, fmt.Errorf("unable to parse body - %s", err.Error())
	}

	identity := &types.DiscordIdentity{}
	err = json.Unmarshal(body, identity)

	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON to struct - %s", err.Error())
	}

	return identity, nil
}

func DiscordVerify(data *types.MakishimaData) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("identity")

		if err != nil {
			data.Logger.Warn().Msg("validation failed - cookie was not found")
			ctx.SetCookie("identity", "", -1, "/", "localhost", true, true)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		claims, err := util.DecodeToken(&cookie)
		if err != nil {
			data.Logger.Warn().Msg(err.Error())
			ctx.SetCookie("identity", "", -1, "/", "localhost", true, true)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		data.Logger.Debug().Msg("validation OK")

		userEntry := database.User{}
		data.Database.Find(&userEntry, claims.DiscordId)
		ctx.JSON(http.StatusOK, gin.H{
			"username": userEntry.Username,
		})
	}
}

// /redirect/discord
func DiscordOAuthRedirect(data *types.MakishimaData) gin.HandlerFunc {
	tokenUrl, _ := url.JoinPath(DiscordEndpoint, "/oauth2/token")
	makishimaId := url.QueryEscape(os.Getenv("MAKISHIMA_ID"))
	makishimaSecret := url.QueryEscape(os.Getenv("MAKISHIMA_SECRET"))
	makishimaRedirect := os.Getenv("MAKISHIMA_REDIRECT")
	makishimaDomain, err := url.Parse(makishimaRedirect)

	if err != nil {
		data.Logger.Fatal().Msgf("redirect URL is invalid - %s", err.Error())
	}

	return func(ctx *gin.Context) {
		code, codeExists := ctx.GetQuery("code")

		if !codeExists {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request",
			})
			return
		}

		// make post request to discord and obtain oauth token
		postForm := url.Values{
			"code":         {url.QueryEscape(code)},
			"grant_type":   {"authorization_code"},
			"redirect_uri": {makishimaRedirect},
		}
		tokenRequest, _ := http.NewRequest("POST", tokenUrl, strings.NewReader(postForm.Encode()))
		tokenRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		tokenRequest.SetBasicAuth(makishimaId, makishimaSecret)
		tokenResponse, err := http.DefaultClient.Do(tokenRequest)

		if err != nil {
			data.Logger.Error().Msgf("Discord OAuth code exchange failed - %s", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "unknown",
			})
			return
		}

		// parse oauth token response
		jsonBytes, err := io.ReadAll(tokenResponse.Body)
		if err != nil {
			data.Logger.Error().Msg("Unable to read token response body")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "unknown",
			})
			return
		}

		if tokenResponse.StatusCode != 200 {
			data.Logger.Error().Msgf("token response returned status code %d:\n%s", tokenResponse.StatusCode, string(jsonBytes))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "unknown",
			})
			return
		}

		//  deserialize token and identify user
		tokenDeserialized := types.DiscordTokenResponse{}
		if json.Unmarshal(jsonBytes, &tokenDeserialized) != nil {
			data.Logger.Error().Msg("unable to serialize token JSON response")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "unknown",
			})
			return
		}

		identity, err := identifyDiscordUser(&tokenDeserialized.AccessToken)
		if err != nil {
			data.Logger.Error().Msg(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "unable to identify user",
			})
			return
		}

		data.Logger.Debug().Msg("successfully identified user")

		// add user to database if ID doesn't exist
		user := database.User{
			ID:           identity.Id,
			Username:     identity.Username,
			AccessToken:  tokenDeserialized.AccessToken,
			RefreshToken: tokenDeserialized.RefreshToken,
			TokenExpiry:  time.Now().Unix() + tokenDeserialized.ExpiresIn,
		}
		dbResult := data.Database.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"access_token"}),
		}).Create(&user)

		if dbResult.Error != nil {
			data.Logger.Error().Msgf("unable to insert to database - %s", dbResult.Error.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "unknown",
			})
			return
		}

		data.Logger.Info().Msgf("inserted into the users table - %d row(s) affected", dbResult.RowsAffected)

		// send back a cookie containing signed jwt
		jwt, err := util.CreateToken(identity, tokenDeserialized.ExpiresIn)
		if err != nil {
			data.Logger.Error().Msgf("unable to craft JWT - %s", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "unknown",
			})
			return
		}

		ctx.SetCookie("identity", jwt, int(tokenDeserialized.ExpiresIn), "/", makishimaDomain.Hostname(), true, true)
		ctx.Redirect(http.StatusTemporaryRedirect, os.Getenv("MAKISHIMA_URI"))
	}
}
