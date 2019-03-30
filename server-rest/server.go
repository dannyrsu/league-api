package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"

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

	results := map[string]interface{}{
		"summonerProfile": summonerProfile,
		"matchHistory":    matchHistory,
	}

	json.NewEncoder(w).Encode(results)
}

func getChampionByKeyHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	champion := models.GetChampionByKey(params.ByName("championkey"))

	json.NewEncoder(w).Encode(champion)
}

func main() {
	router := httprouter.New()
	router.GET("/", defaultHandler)
	router.GET("/v1/summoner/:summonername/stats", getSummonerStatsHandler)
	router.GET("/v1/champion/:championkey", getChampionByKeyHandler)
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
