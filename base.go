package leagueapi

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var riotAPIKey string

func init() {
	envErr := godotenv.Load()

	if envErr != nil {
		log.Fatalf("Error loading .env: %v", envErr)
	}

	riotAPIKey = os.Getenv("RIOT_API_KEY")
}

// GetRiotAPIKey return the key loaded from env
func getRiotAPIKey() string {
	return riotAPIKey
}
