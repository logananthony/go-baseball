package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

// queryPropsTable builds / runs a SELECT * FROM <tbl> with optional filters.
// recognised query params:
//   - any key in exactCols  (exact match)
//   - gamedate (YYYY-MM-DD)  — evaluated in Pacific time, same as /agg-core
//   - limit   (defaults 100, capped 1000)
func (s *APIServer) queryPropsTable(
	w http.ResponseWriter,
	req *http.Request,
	table string,
	exactCols map[string]string, // url param -> column name
) {
	q := req.URL.Query()
	where := []string{}
	args := []interface{}{}
	argID := 1

	// exact-match filters (batterId, pitcherId, gamePk…)
	for param, col := range exactCols {
		if v := q.Get(param); v != "" {
			where = append(where, fmt.Sprintf("%s = $%d", col, argID))
			args = append(args, v)
			argID++
		}
	}

	// Pacific-date filter
	if gd := q.Get("gamedate"); gd != "" {
		// cast through UTC → PT exactly like /agg-core
		where = append(where,
			fmt.Sprintf("(gamedate AT TIME ZONE 'UTC' AT TIME ZONE 'America/Los_Angeles')::date = $%d", argID))
		args = append(args, gd)
		argID++
	}

	sqlParts := []string{fmt.Sprintf("SELECT * FROM %s", table)}
	if len(where) > 0 {
		sqlParts = append(sqlParts, "WHERE "+strings.Join(where, " AND "))
	}

	// limit
	sqlParts = append(sqlParts, "ORDER BY gamedate DESC, gamepk DESC")
	sqlStr := strings.Join(sqlParts, " ")

	rows, err := s.db.Query(sqlStr, args...)
	if err != nil {
		log.Printf("%s query error: %v", table, err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	out := []map[string]interface{}{}

	for rows.Next() {
		data := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range data {
			ptrs[i] = &data[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			log.Printf("scan error: %v", err)
			continue
		}
		row := map[string]interface{}{}
		for i, c := range cols {
			row[c] = *(ptrs[i].(*interface{}))
		}
		out = append(out, row)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func runSims(db *sql.DB, data models.GameData, n int) {
	simData, homeLineup, awayLineup, pitchingSubProbs, homeBullpen, awayBullpen, bullpenRoleProbs, err := sim.PrepareSimData(db, data)
	if err != nil {
		log.Printf("Failed to prepare simulation data: %v", err)
		return
	}
	const maxConcurrentSims = 8 // Tune for your hardware!
	sem := make(chan struct{}, maxConcurrentSims)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			defer func() {
				if r := recover(); r != nil {
					log.Printf("SimulateGame panic: %v", r)
				}
			}()
			sim.SimulateGame(
				db, simData, homeLineup, awayLineup, pitchingSubProbs, homeBullpen, awayBullpen, bullpenRoleProbs, data,
			)
		}()
	}
	wg.Wait()
}

func (s *APIServer) PostSimulateGame(w http.ResponseWriter, req *http.Request) {
	type SimRequest struct {
		UserID string `json:"userId"`
		GamePk int    `json:"gamePk"`
		NSims  int    `json:"nSims"`
	}

	var body SimRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if body.UserID == "" || body.GamePk == 0 || body.NSims == 0 {
		http.Error(w, "Missing one or more required fields", http.StatusBadRequest)
		return
	}

	if body.NSims > 2430 {
		http.Error(w, "Too many simulations requested (max 2430)", http.StatusBadRequest)
		return
	}

	jobID := uuid.New().String()

	gameData := models.GameData{
		GamePk: body.GamePk,
		JobId:  jobID,
	}

	// Create job entry in DB
	_, err := s.db.Exec(`
		INSERT INTO simulation_jobs (id, user_id, status)
		VALUES ($1, $2, 'pending')`, jobID, body.UserID)
	if err != nil {
		log.Printf("Failed to create job in DB: %v", err)
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	// Start simulation in background
	go func(jid, uid string, data models.GameData, n int) {
		_, err := s.db.Exec(`UPDATE simulation_jobs SET status = 'running', updated_at = NOW() WHERE id = $1`, jid)
		if err != nil {
			log.Printf("Failed to update job %s to running: %v", jid, err)
		}

		simData, homeLineup, awayLineup, pitchingSubProbs, homeBullpen, awayBullpen, bullpenRoleProbs, err := sim.PrepareSimData(s.db, data)
		if err != nil {
			log.Printf("Failed to prepare simulation data: %v", err)
			// Mark job as failed and store the error for user feedback:
			_, updErr := s.db.Exec(`
        UPDATE simulation_jobs
        SET status = 'failed', result = $1, updated_at = NOW()
        WHERE id = $2`, err.Error(), jid)
			if updErr != nil {
				log.Printf("Failed to mark job %s as failed: %v", jid, updErr)
			}
			return
		}

		const maxConcurrentSims = 8
		sem := make(chan struct{}, maxConcurrentSims)
		var wg sync.WaitGroup

		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()
				defer func() {
					if r := recover(); r != nil {
						log.Printf("SimulateGame panic: %v", r)
					}
				}()

				sim.SimulateGame(
					s.db, simData, homeLineup, awayLineup, pitchingSubProbs, homeBullpen, awayBullpen, bullpenRoleProbs, data,
				)
			}()
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

	// Respond immediately with job ID
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

// GetAggCore handles GET /agg-core
func (s *APIServer) GetAggCore(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	gamePkStr := q.Get("gamePk")
	limitStr := q.Get("limit")
	gameDateStr := q.Get("gamedate") // e.g., "2025-06-28"

	sqlParts := []string{"SELECT * FROM game_result_agg_core"}
	args := []interface{}{}
	where := []string{}
	argID := 1

	if gamePkStr != "" {
		where = append(where, "gamepk = $"+strconv.Itoa(argID))
		gp, err := strconv.Atoi(gamePkStr)
		if err != nil {
			http.Error(w, "invalid gamePk", http.StatusBadRequest)
			return
		}
		args = append(args, gp)
		argID++
	}

	if gameDateStr != "" {
		where = append(where, "gamedate::date = $"+strconv.Itoa(argID)) // no timezone shifts
		args = append(args, gameDateStr)
		argID++
	}

	if len(where) > 0 {
		sqlParts = append(sqlParts, "WHERE "+strings.Join(where, " AND "))
	}

	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}
	sqlParts = append(sqlParts, fmt.Sprintf("ORDER BY gamepk DESC LIMIT %d", limit))
	sqlStr := strings.Join(sqlParts, " ")

	rows, err := s.db.Query(sqlStr, args...)
	if err != nil {
		log.Printf("agg-core query error: %v", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	results := []map[string]interface{}{}

	for rows.Next() {
		data := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range data {
			ptrs[i] = &data[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			log.Printf("scan error: %v", err)
			continue
		}
		row := map[string]interface{}{}
		for i, c := range cols {
			val := *(ptrs[i].(*interface{}))
			switch v := val.(type) {
			case []byte:
				row[c] = string(v) // decode byte slice to string
			default:
				row[c] = v
			}

		}
		results = append(results, row)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// GET /batter-props
func (s *APIServer) GetBatterProps(w http.ResponseWriter, req *http.Request) {
	s.queryPropsTable(
		w, req,
		"batter_props",
		map[string]string{
			"batterId": "batterid",
			"gamePk":   "gamepk",
		},
	)
}

// GET /pitcher-props
func (s *APIServer) GetPitcherProps(w http.ResponseWriter, req *http.Request) {
	s.queryPropsTable(
		w, req,
		"pitcher_props",
		map[string]string{
			"pitcherId": "pitcherid",
			"gamePk":    "gamepk",
		},
	)
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
	subrouter.HandleFunc("/agg-core", s.GetAggCore).Methods("GET", "OPTIONS")
	subrouter.HandleFunc("/batter-props", s.GetBatterProps).Methods("GET", "OPTIONS")
	subrouter.HandleFunc("/pitcher-props", s.GetPitcherProps).Methods("GET", "OPTIONS")

	log.Printf("Starting API server on %s", s.addr)
	return http.ListenAndServe(s.addr, router)
}
