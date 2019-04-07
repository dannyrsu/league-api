package leagueapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
)

func GetSummonerSpellByKey(spellKey string) map[string]interface{} {
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		panic("no caller information")
	}

	summonerJSON, err := os.Open(path.Dir(filename) + "/static/9.6.1/summoner.json")

	if err != nil {
		log.Fatalf("Error opening summoner.json: %v", err)
	}

	defer summonerJSON.Close()

	sBytes, _ := ioutil.ReadAll(summonerJSON)
	var rawData map[string]interface{}

	jsonErr := json.Unmarshal(sBytes, &rawData)

	if jsonErr != nil {
		log.Fatalf("Error Unmarshaling the data: %v", jsonErr)
	}

	for _, spell := range rawData["data"].(map[string]interface{}) {
		if spell.(map[string]interface{})["key"] == spellKey {
			return spell.(map[string]interface{})
		}
	}

	return nil
}
