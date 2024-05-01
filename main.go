package main

import (
	"log"
	"makishima-backend/routes"
	"makishima-backend/types"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func environCheck(varName string) {
	_, exists := os.LookupEnv(varName)
	if !exists {
		log.Panicf("Environment variable '%s' is empty, terminating early...", varName)
	}
}

func assertEnvironExists() {
	// discord related secrets
	environCheck("MAKISHIMA_ID")
	environCheck("MAKISHIMA_SECRET")
	environCheck("MAKISHIMA_URI")
	environCheck("MAKISHIMA_REDIRECT")
	environCheck("MAKISHIMA_SIGKEY")

	// anilist related secrets
	environCheck("ANILIST_ID")
	environCheck("ANILIST_SECRET")
	environCheck("ANILIST_REDIRECT")

	// database related checks
	environCheck("DATABASE_URI")
}

func main() {
	godotenv.Load()
	assertEnvironExists()

	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()
	db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_URI")), &gorm.Config{})

	if err != nil {
		logger.Fatal().Msgf("unable to connect to database: %s", err.Error())
	}

	data := types.MakishimaData{
		Logger:   &logger,
		Database: db,
	}

	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST"},
		AllowOrigins:     []string{os.Getenv("MAKISHIMA_URI")},
	}))

	// verification
	router.GET("/verify", routes.DiscordVerify(&data))

	// OAuth redirect endpoints
	redirects := router.Group("/redirect")
	redirects.GET("/discord", routes.DiscordOAuthRedirect(&data))
	redirects.GET("/anilist", routes.AnilistOAuthRedirect(&data))

	router.Run("0.0.0.0:3000")
}
