package leagueapi

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

// GetSummonerProfile returns a profile of the summoner
func GetSummonerProfile(summonerName, region string) SummonerProfile {
	apiURL := fmt.Sprintf("https://%v.api.riotgames.com/lol/summoner/v4/summoners/by-name/%v", region, summonerName)

	client := http.Client{
		Timeout: 10 * time.Second,
	}
	req, reqErr := http.NewRequest(http.MethodGet, apiURL, nil)
	req.Header.Set("X-Riot-Token", riotAPIKey)
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

	summoner.MatchHistory = GetMatchSummary(summoner.AccountID, region)

	return summoner
}
