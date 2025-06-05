package api

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/logananthony/go-baseball/pkg/config"
	"github.com/logananthony/go-baseball/pkg/models"
	"github.com/logananthony/go-baseball/pkg/sim"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

// func (data *[]models.GameData) GetSimulateGame(w http.ResponseWriter, req *http.Request) {
func GetSimulateGame(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	homeTeam := query.Get("homeTeam")
	awayTeam := query.Get("awayTeam")
	homeSPStr := query.Get("homeStartingPitcher")
	awaySPStr := query.Get("awayStartingPitcher")
	gameYearStr := query.Get("gameYear")
	nSimsStr := query.Get("nSims")

	// Validate required params
	if homeTeam == "" || awayTeam == "" || homeSPStr == "" || awaySPStr == "" || gameYearStr == "" || nSimsStr == "" {
		http.Error(w, "Missing one or more required query parameters", http.StatusBadRequest)
		return
	}

	// Parse numeric values
	homeSP, err1 := strconv.Atoi(homeSPStr)
	awaySP, err2 := strconv.Atoi(awaySPStr)
	gameYear, err3 := strconv.Atoi(gameYearStr)
	nSims, err4 := strconv.Atoi(nSimsStr)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		http.Error(w, "Invalid query parameter format", http.StatusBadRequest)
		return
	}

	db := config.ConnectDB()
	defer db.Close()

	gameData := models.GameData{
		HomeTeam:            homeTeam,
		AwayTeam:            awayTeam,
		HomeStartingPitcher: homeSP,
		AwayStartingPitcher: awaySP,
		GameYear:            gameYear,
	}

	for i := 0; i < nSims; i++ {
		sim.SimulateGame([]models.GameData{gameData})
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Game simulation complete"))
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1/").Subrouter()
	subrouter.HandleFunc("/simulate", GetSimulateGame).Methods("GET")

	log.Printf("Starting API server on %s", s.addr)
	return http.ListenAndServe(s.addr, router)
}
