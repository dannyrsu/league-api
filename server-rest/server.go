package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dannyrsu/league-api/models"
	"github.com/julienschmidt/httprouter"
)

func defaultHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome to the League of Draaaaven")
}

func getSummonerStatsHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	queryValues := r.URL.Query()

	summonerProfile := models.GetSummonerProfile(params.ByName("summonername"), queryValues.Get("region"))
	matchHistory := models.GetMatchHistory(summonerProfile.AccountID, queryValues.Get("region"), 0, 5)
	realmData := models.GetRealmData()

	results := map[string]interface{}{
		"summonerProfile": summonerProfile,
		"matchHistory":    matchHistory,
		"realmData":       realmData,
	}

	json.NewEncoder(w).Encode(results)
}

func main() {
	router := httprouter.New()
	router.GET("/", defaultHandler)
	router.GET("/v1/summoner/:summonername/stats", getSummonerStatsHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}
