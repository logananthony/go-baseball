package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/logananthony/go-baseball/pkg/models"
	"github.com/logananthony/go-baseball/pkg/sim"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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

func (s *APIServer) GetSimulateGame(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	userID := query.Get("userId")
	homeTeam := query.Get("homeTeam")
	awayTeam := query.Get("awayTeam")
	homeSPStr := query.Get("homeStartingPitcher")
	awaySPStr := query.Get("awayStartingPitcher")
	gameYearStr := query.Get("gameYear")
	nSimsStr := query.Get("nSims")

	// Validate required params
	if userID == "" || homeTeam == "" || awayTeam == "" || homeSPStr == "" || awaySPStr == "" || gameYearStr == "" || nSimsStr == "" {
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

	jobID := uuid.New().String()

	gameData := models.GameData{
		HomeTeam:            homeTeam,
		AwayTeam:            awayTeam,
		HomeStartingPitcher: homeSP,
		AwayStartingPitcher: awaySP,
		GameYear:            gameYear,
	}

	// Insert job into Postgres
	_, err := s.db.Exec(`
		INSERT INTO simulation_jobs (id, user_id, status)
		VALUES ($1, $2, 'pending')`, jobID, userID)
	if err != nil {
		log.Printf("Failed to create job in DB: %v", err)
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	// Launch simulation in background
	go func(jid, uid string, data models.GameData, n int) {
		_, err := s.db.Exec(`UPDATE simulation_jobs SET status = 'running', updated_at = NOW() WHERE id = $1`, jid)
		if err != nil {
			log.Printf("Failed to update job %s to running: %v", jid, err)
		}

		sem := make(chan struct{}, 10)
		var wg sync.WaitGroup

		for i := 0; i < n; i++ {
			sem <- struct{}{}
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				defer func() { <-sem }()
				sim.SimulateGame([]models.GameData{data})
			}(i)
		}

		wg.Wait()

		_, err = s.db.Exec(`
			UPDATE simulation_jobs 
			SET status = 'completed', result = $1, updated_at = NOW()
			WHERE id = $2`, "Simulations complete", jid)
		if err != nil {
			log.Printf("Failed to update job %s to completed: %v", jid, err)
		}
	}(jobID, userID, gameData, nSims)

	// Respond with job ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"jobId": jobID,
	})
}

func (s *APIServer) GetJobStatus(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	userID := query.Get("userId")
	jobID := query.Get("jobId")

	if userID == "" || jobID == "" {
		http.Error(w, "Missing userId or jobId", http.StatusBadRequest)
		return
	}

	var status string
	var result sql.NullString

	err := s.db.QueryRow(`
		SELECT status, result FROM simulation_jobs 
		WHERE id = $1 AND user_id = $2`, jobID, userID).Scan(&status, &result)
	if err == sql.ErrNoRows {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error querying job: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	res := ""
	if result.Valid {
		res = result.String
	}

	json.NewEncoder(w).Encode(map[string]string{
		"jobId":  jobID,
		"status": status,
		"result": res,
	})

}

func (s *APIServer) Run() error {
	router := mux.NewRouter()

	// Add CORS middleware to all routes
	router.Use(corsMiddleware)

	subrouter := router.PathPrefix("/api/v1/").Subrouter()
	subrouter.HandleFunc("/simulate", s.GetSimulateGame).Methods("GET", "OPTIONS")
	subrouter.HandleFunc("/status", s.GetJobStatus).Methods("GET", "OPTIONS")

	log.Printf("Starting API server on %s", s.addr)
	return http.ListenAndServe(s.addr, router)
}
