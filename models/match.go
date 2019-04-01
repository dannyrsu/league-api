package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// MatchHistory model
type MatchHistory struct {
	Matches []struct {
		Lane       string `json:"lane"`
		GameID     int64  `json:"gameId"`
		Champion   int32  `json:"champion"`
		PlatformID string `json:"platformId"`
		Timestamp  int64  `json:"timestamp"`
		Queue      int32  `json:"queue"`
		Role       string `json:"role"`
		Season     int32  `json:"season"`
	} `json:"matches"`
	EndIndex   int32 `json:"endIndex"`
	StartIndex int32 `json:"startIndex"`
	TotalGames int32 `json:"totalGames"`
}

// GetMatchHistory for the account
func GetMatchHistory(accountID, region string, beginIndex, endIndex int) MatchHistory {
	apiURL := fmt.Sprintf("https://%v.api.riotgames.com/lol/match/v4/matchlists/by-account/%v?queue=420&endIndex=%v&beginIndex=%v", region, accountID, endIndex, beginIndex)
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	req, reqErr := http.NewRequest(http.MethodGet, apiURL, nil)
	req.Header.Set("X-Riot-Token", GetRiotAPIKey())
	if reqErr != nil {
		log.Fatalf("Error creating request: %v", reqErr)
	}
	resp, getErr := client.Do(req)

	if getErr != nil {
		log.Fatalf("Error getting summoner matches: %v", getErr)
	}

	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatalf("Error reading response body: %v", readErr)
	}

	matchHistory := MatchHistory{}

	jsonErr := json.Unmarshal(body, &matchHistory)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return matchHistory
}

// GetGameData return data for a game id
func GetGameData(matchID int64, region string) map[string]interface{} {
	apiURL := fmt.Sprintf("https://%v.api.riotgames.com/lol/match/v4/matches/%v", region, matchID)

	client := http.Client{
		Timeout: 10 * time.Second,
	}
	req, reqErr := http.NewRequest(http.MethodGet, apiURL, nil)
	req.Header.Set("X-Riot-Token", GetRiotAPIKey())
	if reqErr != nil {
		log.Fatalf("Error creating request: %v", reqErr)
	}
	resp, getErr := client.Do(req)

	if getErr != nil {
		log.Fatalf("Error getting match data: %v", getErr)
	}

	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatalf("Error reading response body: %v", readErr)
	}

	var result map[string]interface{}

	jsonErr := json.Unmarshal(body, &result)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return result
}
