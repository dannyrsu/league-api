package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func GetRealmData() map[string]interface{} {
	realmData := make(map[string]interface{})
	realmURL := "https://ddragon.leagueoflegends.com/realms/na.json"

	res, getErr := http.Get(realmURL)

	if getErr != nil {
		log.Fatalf("Error getting the realm file: %v", getErr)
	}
	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)

	if readErr != nil {
		log.Fatalf("Error reading realm file: %v", readErr)
	}

	jsonErr := json.Unmarshal(body, &realmData)

	if jsonErr != nil {
		log.Fatalf("Error unmarshaling json data: %v", jsonErr)
	}

	return realmData
}
