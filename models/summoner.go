package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// SummonerProfile model
type SummonerProfile struct {
	ProfileIconID int32  `json:"profileIconId"`
	Name          string `json:"name"`
	PUUID         string `json:"puuid"`
	SummonerLevel int64  `json:"summonerLevel"`
	RevisionDate  int64  `json:"revisionDate"`
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
}

// MatchHistory model
type MatchHistory struct {
	Matches []struct {
		Lane       string                 `json:"lane"`
		GameID     int64                  `json:"gameId"`
		Champion   int32                  `json:"champion"`
		PlatformID string                 `json:"platformId"`
		Timestamp  int64                  `json:"timestamp"`
		Queue      int32                  `json:"queue"`
		Role       string                 `json:"role"`
		Season     int32                  `json:"season"`
		Game       map[string]interface{} `json:"game"`
	} `json:"matches"`
	EndIndex   int32 `json:"endIndex"`
	StartIndex int32 `json:"startIndex"`
	TotalGames int32 `json:"totalGames"`
}

// GetSummonerProfile returns a profile of the summoner
func GetSummonerProfile(summonerName, region string) SummonerProfile {
	apiURL := fmt.Sprintf("https://%v.api.riotgames.com/lol/summoner/v4/summoners/by-name/%v", region, summonerName)

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

	for i := range matchHistory.Matches {
		matchHistory.Matches[i].Game = GetGameData(matchHistory.Matches[i].GameID, region)
	}

	return matchHistory
}
