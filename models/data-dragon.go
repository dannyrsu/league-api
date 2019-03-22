package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type GameVersions struct {
	Item        string `json:"item"`
	Rune        string `json:"rune"`
	Mastery     string `json:"mastery"`
	Summoner    string `json:"summoner"`
	Champion    string `json:"champion"`
	ProfileIcon string `json:"profileicon"`
	Map         string `json:"map"`
	Language    string `json:"language"`
	Sticker     string `json:"sticker"`
}

type RealmData struct {
	GameVersions   GameVersions `json:"n"`
	Version        string       `json:"v"`
	Locale         string       `json:"l"`
	CDN            string       `json:"cdn"`
	DataDragon     string       `json:"dd"`
	LG             string       `json:"lg"`
	CSS            string       `json:"css"`
	ProfileIconMax int          `json:"profileiconmax"`
	Store          string       `json:"-"`
}

func getRealmData() RealmData {
	realmData := RealmData{}
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
