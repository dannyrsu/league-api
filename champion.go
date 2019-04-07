package leagueapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
)

// GetChampionByKey return champion data from static files
func GetChampionByKey(championKey string) map[string]interface{} {
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		panic("no caller information")
	}

	championJSON, err := os.Open(path.Dir(filename) + "/static/9.6.1/championFull.json")
	if err != nil {
		log.Fatalf("Error opening champion.json: %v", err)
	}

	defer championJSON.Close()

	cBytes, _ := ioutil.ReadAll(championJSON)
	var rawData map[string]interface{}

	jsonErr := json.Unmarshal(cBytes, &rawData)
	if jsonErr != nil {
		log.Fatalf("Error Unmarshaling the data: %v", jsonErr)
	}

	for _, champion := range rawData["data"].(map[string]interface{}) {
		if champion.(map[string]interface{})["key"] == championKey {
			return champion.(map[string]interface{})
		}
	}

	return nil
}
