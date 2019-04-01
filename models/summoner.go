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
	ProfileIconID int32                    `json:"profileIconId"`
	Name          string                   `json:"name"`
	PUUID         string                   `json:"puuid"`
	SummonerLevel int64                    `json:"summonerLevel"`
	RevisionDate  int64                    `json:"revisionDate"`
	ID            string                   `json:"id"`
	AccountID     string                   `json:"accountId"`
	MatchHistory  []map[string]interface{} `json:"matchHistory"`
}

func getMatchSummary(accountID, region string) []map[string]interface{} {
	matchHistory := GetMatchHistory(accountID, region, 0, 5)
	matchSummary := make([]map[string]interface{}, len(matchHistory.Matches))

	for i := 0; i < len(matchHistory.Matches); i++ {
		match := matchHistory.Matches[i]
		game := GetGameData(match.GameID, region)
		var participantID float64
		for j := 0; j < len(game["participantIdentities"].([]interface{})); j++ {
			participant := game["participantIdentities"].([]interface{})[j].(map[string]interface{})
			player := participant["player"]
			if player.(map[string]interface{})["accountId"] == accountID {
				participantID = participant["participantId"].(float64)
				break
			}
		}

		if participantID > 0 {
			for k := 0; k < len(game["participants"].([]interface{})); k++ {
				participant := game["participants"].([]interface{})[k].(map[string]interface{})
				if participant["participantId"] == participantID {
					matchSummary[i] = participant
					break
				}
			}
		}
	}

	return matchSummary
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

	summoner.MatchHistory = getMatchSummary(summoner.AccountID, region)

	return summoner
}
