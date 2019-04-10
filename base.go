package leagueapi

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var riotAPIKey, staticFilesRoot string

func init() {
	envErr := godotenv.Load()

	if envErr != nil {
		log.Fatalf("Error loading .env: %v", envErr)
	}

	riotAPIKey = os.Getenv("RIOT_API_KEY")
	staticFilesRoot = os.Getenv("RIOT_STATIC_FILES_ROOT")
}
