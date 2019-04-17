package leagueapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/mitchellh/mapstructure"
)

type SummonerSpell struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Tooltip      string `json:"tooltip"`
	Maxrank      int    `json:"maxrank"`
	Cooldown     []int  `json:"cooldown"`
	CooldownBurn string `json:"cooldownBurn"`
	Cost         []int  `json:"cost"`
	CostBurn     string `json:"costBurn"`
	Datavalues   struct {
	} `json:"datavalues"`
	Effect        []interface{} `json:"effect"`
	EffectBurn    []interface{} `json:"effectBurn"`
	Vars          []interface{} `json:"vars"`
	Key           string        `json:"key"`
	SummonerLevel int           `json:"summonerLevel"`
	Modes         []string      `json:"modes"`
	CostType      string        `json:"costType"`
	Maxammo       string        `json:"maxammo"`
	Range         []int         `json:"range"`
	RangeBurn     string        `json:"rangeBurn"`
	Image         struct {
		Full   string `json:"full"`
		Sprite string `json:"sprite"`
		Group  string `json:"group"`
		X      int    `json:"x"`
		Y      int    `json:"y"`
		W      int    `json:"w"`
		H      int    `json:"h"`
	} `json:"image"`
	Resource string `json:"resource"`
}

func GetSummonerSpellByKey(spellKey string) SummonerSpell {
	summonerJSON, err := os.Open(staticFilesRoot + "summoner.json")
	var summonerSpell SummonerSpell

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
			errDecode := mapstructure.Decode(spell, &summonerSpell)

			if errDecode != nil {
				log.Fatalf("Error converting map to struct: %v", errDecode)
			}

			break
		}
	}

	return summonerSpell
}
