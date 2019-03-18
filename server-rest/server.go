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
	summonerProfile := models.GetSummonerProfile(params.ByName("summonername"), params.ByName("region"))
	json.NewEncoder(w).Encode(summonerProfile)
}

func main() {
	router := httprouter.New()
	router.GET("/", defaultHandler)
	router.GET("/summoner/:summonername/region/:region/stats", getSummonerStatsHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}
