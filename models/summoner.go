package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// SummonerProfile model
type SummonerProfile struct {
	ProfileIconID int    `json:"profileIconId"`
	Name          string `json:"name"`
	PUUID         string `json:"puuid"`
	SummonerLevel int    `json:"summonerLevel"`
	RevisionDate  int    `json:"revisionDate"`
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
}

// GetSummonerProfile returns a profile of the summoner
func GetSummonerProfile(apiURL string) SummonerProfile {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	req, reqErr := http.NewRequest(http.MethodGet, apiURL, nil)
	if reqErr != nil {
		log.Fatalf("Error creating request: %v", reqErr)
	}
	resp, getErr := client.Do(req)

	if getErr != nil {
		log.Fatalf("Error getting summoner profile: %v", getErr)
	}

	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatalf("Error reading response body: %v", readErr)
	}

	summoner := SummonerProfile{}

	jsonErr := json.Unmarshal(body, &summoner)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return summoner
}