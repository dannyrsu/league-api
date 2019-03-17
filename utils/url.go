package utils

import "fmt"

// GetSummonerProfileURL construct url TODO:there probably is a better way to do this....
func GetSummonerProfileURL(summonerName, region, apiKey string) string {
	return fmt.Sprintf("https://%v.api.riotgames.com/lol/summoner/v4/summoners/by-name/%v?api_key=%v", region, summonerName, apiKey)
}
