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

type server struct {
	router *httprouter.Router
}

func (*server) defaultHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome to the League of Draaaaven")
}

func (*server) getSummonerStatsHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	queryValues := r.URL.Query()

	summonerProfile := models.GetSummonerProfile(params.ByName("summonername"), queryValues.Get("region"))

	results := map[string]interface{}{
		"summonerProfile": summonerProfile,
	}

	json.NewEncoder(w).Encode(results)
}

func (*server) getChampionByKeyHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	champion := models.GetChampionByKey(params.ByName("championkey"))

	json.NewEncoder(w).Encode(champion)
}

func (s *server) routes() {
	s.router.GET("/", s.defaultHandler)
	s.router.GET("/v1/summoner/:summonername/stats", s.getSummonerStatsHandler)
	s.router.GET("/v1/champion/:championkey", s.getChampionByKeyHandler)
	s.router.ServeFiles("/static/*filepath", http.Dir("static"))
}

func main() {
	server := &server{
		router: httprouter.New(),
	}

	server.routes()

	handler := cors.Default().Handler(server.router)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
