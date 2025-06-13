package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func (s *APIServer) PostSimulateGame(w http.ResponseWriter, req *http.Request) {
	type SimRequest struct {
		UserID              string `json:"userId"`
		HomeTeam            string `json:"homeTeam"`
		AwayTeam            string `json:"awayTeam"`
		HomeStartingPitcher int    `json:"homeStartingPitcher"`
		AwayStartingPitcher int    `json:"awayStartingPitcher"`
		GameYear            int    `json:"gameYear"`
		NSims               int    `json:"nSims"`
	}

	var body SimRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if body.UserID == "" || body.HomeTeam == "" || body.AwayTeam == "" ||
		body.HomeStartingPitcher == 0 || body.AwayStartingPitcher == 0 ||
		body.GameYear == 0 || body.NSims == 0 {
		http.Error(w, "Missing one or more required fields", http.StatusBadRequest)
		return
	}

	jobID := uuid.New().String()

	gameData := models.GameData{
		HomeTeam:            body.HomeTeam,
		AwayTeam:            body.AwayTeam,
		HomeStartingPitcher: body.HomeStartingPitcher,
		AwayStartingPitcher: body.AwayStartingPitcher,
		GameYear:            body.GameYear,
		JobId:               jobID,
	}

	_, err := s.db.Exec(`
		INSERT INTO simulation_jobs (id, user_id, status)
		VALUES ($1, $2, 'pending')`, jobID, body.UserID)
	if err != nil {
		log.Printf("Failed to create job in DB: %v", err)
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

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
	}(jobID, body.UserID, gameData, body.NSims)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"jobId": jobID,
	})
}

func (s *APIServer) PostJobStatus(w http.ResponseWriter, req *http.Request) {
	type StatusRequest struct {
		UserID string `json:"userId"`
		JobID  string `json:"jobId"`
	}

	var body StatusRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if body.UserID == "" || body.JobID == "" {
		http.Error(w, "Missing userId or jobId", http.StatusBadRequest)
		return
	}

	var status string
	var result sql.NullString

	err := s.db.QueryRow(`
		SELECT status, result FROM simulation_jobs 
		WHERE id = $1 AND user_id = $2`, body.JobID, body.UserID).Scan(&status, &result)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"jobId":  body.JobID,
		"status": status,
		"result": res,
	})
}

func (s *APIServer) GetGameResults(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	where := []string{}
	args := []interface{}{}
	argID := 1

	// Exact-match filters
	exactFilters := map[string]string{
		"velocity":       "velocity",
		"exit_velocity":  "exit_velocity",
		"launch_angle":   "launch_angle",
		"spray_angle":    "spray_angle",
		"pitch_type":     "pitch_type",
		"batterid":       "batterid",
		"pitcherid":      "pitcherid",
		"game_year":      "game_year",
		"zone":           "zone",
		"batter_stands":  "batter_stands",
		"event_type":     "event_type",
		"jobid":          "jobid",
		"inning":         "inning",
		"inning_topbot":  "inning_topbot",
		"pitcher_throws": "pitcher_throws",
		"on1b":           "on1b",
		"on2b":           "on2b",
		"on3b":           "on3b",
	}

	// Range filters: field -> SQL clause
	rangeFilters := map[string]string{
		"velocityMin":     "velocity >= ",
		"velocityMax":     "velocity <= ",
		"exitVelocityMin": "exit_velocity >= ",
		"exitVelocityMax": "exit_velocity <= ",
		"launchAngleMin":  "launch_angle >= ",
		"launchAngleMax":  "launch_angle <= ",
		"sprayAngleMin":   "spray_angle >= ",
		"sprayAngleMax":   "spray_angle <= ",
		"inningMin":       "inning >= ",
		"inningMax":       "inning <= ",
	}

	// Add exact match conditions
	for param, column := range exactFilters {
		if val := query.Get(param); val != "" {
			where = append(where, column+" = $"+strconv.Itoa(argID))
			args = append(args, val)
			argID++
		}
	}

	// Add range conditions
	for param, clause := range rangeFilters {
		if val := query.Get(param); val != "" {
			where = append(where, clause+"$"+strconv.Itoa(argID))
			args = append(args, val)
			argID++
		}
	}

	sqlStr := "SELECT * FROM game_result"
	if len(where) > 0 {
		sqlStr += " WHERE " + strings.Join(where, " AND ")
	}
	sqlStr += " ORDER BY created_at DESC LIMIT 100" // safety limit

	rows, err := s.db.Query(sqlStr, args...)
	if err != nil {
		log.Printf("DB query error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	results := []map[string]interface{}{}

	cols, _ := rows.Columns()
	for rows.Next() {
		colsData := make([]interface{}, len(cols))
		colsPtrs := make([]interface{}, len(cols))
		for i := range colsData {
			colsPtrs[i] = &colsData[i]
		}

		if err := rows.Scan(colsPtrs...); err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}

		row := make(map[string]interface{})
		for i, colName := range cols {
			val := colsPtrs[i].(*interface{})
			row[colName] = *val
		}
		results = append(results, row)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (s *APIServer) GetJobStatusQueryParams(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	userID := query.Get("userId")
	jobID := query.Get("jobId")

	if userID == "" {
		http.Error(w, "Missing userId", http.StatusBadRequest)
		return
	}

	// Case 1: Specific job requested
	if jobID != "" {
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

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"jobId":  jobID,
			"status": status,
			"result": res,
		})
		return
	}

	// Case 2: Return all jobs for user
	rows, err := s.db.Query(`
		SELECT id, status, result FROM simulation_jobs 
		WHERE user_id = $1 ORDER BY updated_at DESC`, userID)
	if err != nil {
		log.Printf("Error querying jobs for user %s: %v", userID, err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var jobs []map[string]string

	for rows.Next() {
		var id, status string
		var result sql.NullString

		if err := rows.Scan(&id, &status, &result); err != nil {
			log.Printf("Error scanning job row: %v", err)
			continue
		}

		res := ""
		if result.Valid {
			res = result.String
		}

		jobs = append(jobs, map[string]string{
			"jobId":  id,
			"status": status,
			"result": res,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()

	// Add CORS middleware to all routes
	router.Use(corsMiddleware)

	subrouter := router.PathPrefix("/api/v1/").Subrouter()
	subrouter.HandleFunc("/simulate", s.PostSimulateGame).Methods("POST", "OPTIONS")
	subrouter.HandleFunc("/status", s.PostJobStatus).Methods("POST", "OPTIONS")
	subrouter.HandleFunc("/results", s.GetGameResults).Methods("GET", "OPTIONS")
	subrouter.HandleFunc("/status", s.GetJobStatusQueryParams).Methods("GET", "OPTIONS")

	log.Printf("Starting API server on %s", s.addr)
	return http.ListenAndServe(s.addr, router)
}
