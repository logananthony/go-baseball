package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/logananthony/go-baseball/pkg/models"
	"github.com/logananthony/go-baseball/pkg/poster"
	"github.com/logananthony/go-baseball/pkg/sim"
	"github.com/logananthony/go-baseball/pkg/utils"
)

func corsMiddleware(next http.Handler) http.Handler {
	allowed := map[string]bool{
		"http://localhost:3000":       true,
		"https://go-baseball.com":     true,
		"https://www.go-baseball.com": true,
		// add any other frontends here
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && allowed[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			// If you use cookies/session:
			// w.Header().Set("Access-Control-Allow-Credentials", "true")
			// If the client needs to read custom response headers:
			// w.Header().Set("Access-Control-Expose-Headers", "X-Custom-Header")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent) // 204
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
	// Ensure new progress columns exist
	if _, err := db.Exec(`ALTER TABLE simulation_jobs ADD COLUMN IF NOT EXISTS total_simulations INT DEFAULT 0`); err != nil {
		log.Printf("failed adding total_simulations column: %v", err)
	}
	if _, err := db.Exec(`ALTER TABLE simulation_jobs ADD COLUMN IF NOT EXISTS current_simulation INT DEFAULT 0`); err != nil {
		log.Printf("failed adding current_simulation column: %v", err)
	}
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
		startUTC, endUTC, err := pacificDayRangeUTC(gd)
		if err != nil {
			http.Error(w, "invalid gamedate", http.StatusBadRequest)
			return
		}
		where = append(where, fmt.Sprintf("gamedate >= $%d AND gamedate < $%d", argID, argID+1))
		args = append(args, startUTC, endUTC)
		argID += 2
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

func runSims(db *sql.DB, data models.GameData, n int, outcomeMode string) {
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
				db, simData, homeLineup, awayLineup, pitchingSubProbs, homeBullpen, awayBullpen, bullpenRoleProbs, data, outcomeMode,
			)
		}()
	}
	wg.Wait()
}

func aggregateEventCounts(eventCounts []int) (prob05, prob15, avg, total, iqr, q80, lower05, upper05, lower15, upper15 float64) {
	n := len(eventCounts)
	if n == 0 {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
	}
	counts := make([]float64, n)
	for i, c := range eventCounts {
		counts[i] = float64(c)
	}
	avg = utils.Mean(counts)
	total = float64(len(eventCounts))
	iqr = utils.QuantileWidth(counts, 0.25, 0.75)
	q80 = utils.QuantileWidth(counts, 0.10, 0.90)

	num05, num15 := 0, 0
	for _, c := range eventCounts {
		if c >= 1 {
			num05++
		}
		if c >= 2 {
			num15++
		}
	}
	prob05 = float64(num05) / float64(n)
	prob15 = float64(num15) / float64(n)

	lower05, upper05 = utils.BinomialCI(prob05, n)
	lower15, upper15 = utils.BinomialCI(prob15, n)
	return
}

func getKProp(kCounts []int, threshold int) (prob, lower95, upper95 float64) {
	n, numOver := len(kCounts), 0
	if n == 0 {
		return 0, 0, 0
	}
	for _, v := range kCounts {
		if v >= threshold {
			numOver++
		}
	}
	prob = float64(numOver) / float64(n)
	lower95, upper95 = utils.BinomialCI(prob, n)
	return
}

func (s *APIServer) PostSimulateGame(w http.ResponseWriter, req *http.Request) {
	type SimRequest struct {
		UserID      string `json:"userId"`
		GamePk      int    `json:"gamePk"`
		NSims       int    `json:"nSims"`
		OutcomeMode string `json:"outcomeMode"`
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
               INSERT INTO simulation_jobs (id, user_id, status, total_simulations, current_simulation)
               VALUES ($1, $2, 'pending', $3, 0)`, jobID, body.UserID, body.NSims)
	if err != nil {
		log.Printf("Failed to create job in DB: %v", err)
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	// Start simulation in background
	go func(jid, uid string, data models.GameData, n int, outcomeMode string) {
		_, err := s.db.Exec(`UPDATE simulation_jobs SET status = 'running', updated_at = NOW() WHERE id = $1`, jid)
		if err != nil {
			log.Printf("Failed to update job %s to running: %v", jid, err)
		}

		simData, homeLineup, awayLineup, pitchingSubProbs, homeBullpen, awayBullpen, bullpenRoleProbs, err := sim.PrepareSimData(s.db, data)
		// ----- Add this block right after getting the lineups -----

		// Build batter metadata (batterID -> name, order, team)
		batterMeta := map[int64]struct {
			Name  string
			Order int
			Team  string
		}{}

		pitcherMeta := map[int64]struct {
			Name string
			Team string
		}{}
		// Build pitcher metadata (pitcherID -> name)
		// pitcherMeta := map[int64]string{}

		// Loop through BOTH lineups (home and away)
		for _, b := range append(homeLineup, awayLineup...) {
			// Batter
			batterMeta[int64(b.PlayerId)] = struct {
				Name  string
				Order int
				Team  string
			}{
				Name:  b.PlayerName,
				Order: b.BattingOrder,
				Team:  b.Team,
			}
			if b.Team == "home" {
				if _, ok := pitcherMeta[int64(b.HomePitcherId)]; !ok {
					pitcherMeta[int64(b.HomePitcherId)] = struct{ Name, Team string }{
						Name: b.HomePitcherName, Team: "home",
					}
				}
			} else {
				if _, ok := pitcherMeta[int64(b.AwayPitcherId)]; !ok {
					pitcherMeta[int64(b.AwayPitcherId)] = struct{ Name, Team string }{
						Name: b.AwayPitcherName, Team: "away",
					}
				}
			}
		}

		playerByID := make(map[int]models.MLBPlayerInfo, len(simData.PlayerInfo))
		for i := range simData.PlayerInfo {
			if simData.PlayerInfo[i].ID != nil {
				playerByID[*simData.PlayerInfo[i].ID] = simData.PlayerInfo[i]
			}
		}

		// pitcherMeta := map[int64]struct{ Name, Team string }{}

		for _, side := range []struct {
			o     *models.BullpenOrder
			label string
		}{
			{homeBullpen, "home"},
			{awayBullpen, "away"},
		} {
			if side.o == nil {
				continue
			}
			for _, id := range []int{side.o.PlayerID1, side.o.PlayerID2, side.o.PlayerID3, side.o.PlayerID4, side.o.PlayerID5, side.o.PlayerID6, side.o.PlayerID7, side.o.PlayerID8} {
				if id == 0 {
					continue
				}
				if p, ok := playerByID[id]; ok {
					name := ""
					if p.FullName != nil && *p.FullName != "" {
						name = *p.FullName
					} else {
						if p.FirstName != nil {
							name += *p.FirstName
						}
						if p.LastName != nil {
							if name != "" {
								name += " "
							}
							name += *p.LastName
						}
					}
					pitcherMeta[int64(id)] = struct{ Name, Team string }{name, side.label} // "home"/"away"
				}
			}
		}

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

		maxWorkers := runtime.NumCPU()
		if maxWorkers > 4 {
			maxWorkers = 4
		}
		if n < maxWorkers {
			maxWorkers = n
		}

		jobs := make(chan int)
		var wg sync.WaitGroup

		var (
			homeTeam, awayTeam, homePitcherName, awayPitcherName string
			gameDate                                             time.Time
			metaOnce                                             sync.Once
		)

		outcomes := make([]models.SimOutcome, n)
		var completed int64

		worker := func() {
			defer wg.Done()
			for i := range jobs {
				attempt := 0
				for {
					var err error
					func() {
						defer func() {
							if r := recover(); r != nil {
								err = fmt.Errorf("SimulateGame panic: %v", r)
							}
						}()
						h, a, ht, at, hpn, apn, gd, batterEvents, pitcherEvents, runEvents := sim.SimulateGame(
							s.db, simData, homeLineup, awayLineup, pitchingSubProbs, homeBullpen, awayBullpen, bullpenRoleProbs, data, outcomeMode,
						)
						metaOnce.Do(func() {
							homeTeam = ht
							awayTeam = at
							homePitcherName = hpn
							awayPitcherName = apn
							gameDate = gd
						})
						outcomes[i] = models.SimOutcome{
							HomeScore:     h,
							AwayScore:     a,
							BatterEvents:  batterEvents,
							PitcherEvents: pitcherEvents,
							RunEvents:     runEvents,
						}
					}()
					if err == nil {
						newVal := atomic.AddInt64(&completed, 1)
						if _, dbErr := s.db.Exec(`UPDATE simulation_jobs SET current_simulation = $1, updated_at = NOW() WHERE id = $2`, newVal, jid); dbErr != nil {
							log.Printf("Failed to update progress for job %s: %v", jid, dbErr)
						}
						break
					}
					attempt++
					if attempt >= 3 {
						log.Printf("SimulateGame failed after retries: %v", err)
						break
					}
					time.Sleep(time.Duration(1<<attempt) * time.Second)
				}
			}
		}

		for w := 0; w < maxWorkers; w++ {
			wg.Add(1)
			go worker()
		}

		batchSize := maxWorkers
		for start := 0; start < n; start += batchSize {
			end := start + batchSize
			if end > n {
				end = n
			}
			for i := start; i < end; i++ {
				jobs <- i
			}
			time.Sleep(200 * time.Millisecond)
		}
		close(jobs)
		wg.Wait()

		// ========== GAME-LEVEL AGGREGATION ==========

		// --- Helper functions ---
		gameTotalProb := func(vals []int, threshold int) (float64, float64, float64) {
			n := len(vals)
			if n == 0 {
				return 0, 0, 0
			}
			cnt := 0
			for _, x := range vals {
				if x > threshold {
					cnt++
				}
			}
			prob := float64(cnt) / float64(n)
			lower, upper := utils.BinomialCI(prob, n)
			return prob, lower, upper
		}
		gameTotalProbIncl := func(vals []int, threshold int) (float64, float64, float64) {
			n := len(vals)
			if n == 0 {
				return 0, 0, 0
			}
			cnt := 0
			for _, x := range vals {
				if x >= threshold {
					cnt++
				}
			}
			prob := float64(cnt) / float64(n)
			lower, upper := utils.BinomialCI(prob, n)
			return prob, lower, upper
		}

		// --- Build slices ---
		var (
			totalRuns  []int
			homeScores []int
			awayScores []int
			spreads    []int // homeScore - awayScore
		)
		for _, sim := range outcomes {
			t := sim.HomeScore + sim.AwayScore
			totalRuns = append(totalRuns, t)
			homeScores = append(homeScores, sim.HomeScore)
			awayScores = append(awayScores, sim.AwayScore)
			spreads = append(spreads, sim.HomeScore-sim.AwayScore)
		}

		loc, _ := time.LoadLocation("America/Los_Angeles")
		gameDatePacific := gameDate.In(loc)

		agg := models.GameResultAggCore{
			GamePk:          int64(data.GamePk),
			HomeTeamAbbr:    homeTeam,
			AwayTeamAbbr:    awayTeam,
			HomePitcherName: homePitcherName,
			AwayPitcherName: awayPitcherName,
			GameDate:        gameDatePacific,
		}

		// Totals (over X.5)
		agg.TotalOver15, agg.TotalOver15Lower95, agg.TotalOver15Upper95 = gameTotalProb(totalRuns, 1)
		agg.TotalOver25, agg.TotalOver25Lower95, agg.TotalOver25Upper95 = gameTotalProb(totalRuns, 2)
		agg.TotalOver35, agg.TotalOver35Lower95, agg.TotalOver35Upper95 = gameTotalProb(totalRuns, 3)
		agg.TotalOver45, agg.TotalOver45Lower95, agg.TotalOver45Upper95 = gameTotalProb(totalRuns, 4)
		agg.TotalOver55, agg.TotalOver55Lower95, agg.TotalOver55Upper95 = gameTotalProb(totalRuns, 5)
		agg.TotalOver65, agg.TotalOver65Lower95, agg.TotalOver65Upper95 = gameTotalProb(totalRuns, 6)
		agg.TotalOver75, agg.TotalOver75Lower95, agg.TotalOver75Upper95 = gameTotalProb(totalRuns, 7)
		agg.TotalOver85, agg.TotalOver85Lower95, agg.TotalOver85Upper95 = gameTotalProb(totalRuns, 8)
		agg.TotalOver95, agg.TotalOver95Lower95, agg.TotalOver95Upper95 = gameTotalProb(totalRuns, 9)
		agg.TotalOver105, agg.TotalOver105Lower95, agg.TotalOver105Upper95 = gameTotalProb(totalRuns, 10)
		agg.TotalOver115, agg.TotalOver115Lower95, agg.TotalOver115Upper95 = gameTotalProb(totalRuns, 11)
		agg.TotalOver125, agg.TotalOver125Lower95, agg.TotalOver125Upper95 = gameTotalProb(totalRuns, 12)

		// Home/Away team totals (over X.5)
		agg.HomeTotalOver05, agg.HomeTotalOver05Lower95, agg.HomeTotalOver05Upper95 = gameTotalProbIncl(homeScores, 1)
		agg.HomeTotalOver15, agg.HomeTotalOver15Lower95, agg.HomeTotalOver15Upper95 = gameTotalProb(homeScores, 1)
		agg.HomeTotalOver25, agg.HomeTotalOver25Lower95, agg.HomeTotalOver25Upper95 = gameTotalProb(homeScores, 2)
		agg.HomeTotalOver35, agg.HomeTotalOver35Lower95, agg.HomeTotalOver35Upper95 = gameTotalProb(homeScores, 3)
		agg.HomeTotalOver45, agg.HomeTotalOver45Lower95, agg.HomeTotalOver45Upper95 = gameTotalProb(homeScores, 4)
		agg.HomeTotalOver55, agg.HomeTotalOver55Lower95, agg.HomeTotalOver55Upper95 = gameTotalProb(homeScores, 5)
		agg.HomeTotalOver65, agg.HomeTotalOver65Lower95, agg.HomeTotalOver65Upper95 = gameTotalProb(homeScores, 6)

		agg.AwayTotalOver05, agg.AwayTotalOver05Lower95, agg.AwayTotalOver05Upper95 = gameTotalProbIncl(awayScores, 1)
		agg.AwayTotalOver15, agg.AwayTotalOver15Lower95, agg.AwayTotalOver15Upper95 = gameTotalProb(awayScores, 1)
		agg.AwayTotalOver25, agg.AwayTotalOver25Lower95, agg.AwayTotalOver25Upper95 = gameTotalProb(awayScores, 2)
		agg.AwayTotalOver35, agg.AwayTotalOver35Lower95, agg.AwayTotalOver35Upper95 = gameTotalProb(awayScores, 3)
		agg.AwayTotalOver45, agg.AwayTotalOver45Lower95, agg.AwayTotalOver45Upper95 = gameTotalProb(awayScores, 4)
		agg.AwayTotalOver55, agg.AwayTotalOver55Lower95, agg.AwayTotalOver55Upper95 = gameTotalProb(awayScores, 5)
		agg.AwayTotalOver65, agg.AwayTotalOver65Lower95, agg.AwayTotalOver65Upper95 = gameTotalProb(awayScores, 6)

		// Spread: (home - away) > X.5 (ex: minus_15 = home wins by 2+)
		spreadProb := func(vals []int, thresh float64) (float64, float64, float64) {
			n := len(vals)
			if n == 0 {
				return 0, 0, 0
			}
			cnt := 0
			for _, x := range vals {
				if float64(x) > thresh {
					cnt++
				}
			}
			prob := float64(cnt) / float64(n)
			lower, upper := utils.BinomialCI(prob, n)
			return prob, lower, upper
		}
		agg.SpreadMinus55, agg.SpreadMinus55Lower95, agg.SpreadMinus55Upper95 = spreadProb(spreads, 5.5)
		agg.SpreadMinus45, agg.SpreadMinus45Lower95, agg.SpreadMinus45Upper95 = spreadProb(spreads, 4.5)
		agg.SpreadMinus35, agg.SpreadMinus35Lower95, agg.SpreadMinus35Upper95 = spreadProb(spreads, 3.5)
		agg.SpreadMinus25, agg.SpreadMinus25Lower95, agg.SpreadMinus25Upper95 = spreadProb(spreads, 2.5)
		agg.SpreadMinus15, agg.SpreadMinus15Lower95, agg.SpreadMinus15Upper95 = spreadProb(spreads, 1.5)
		agg.SpreadPlus15, agg.SpreadPlus15Lower95, agg.SpreadPlus15Upper95 = spreadProb(spreads, -1.5)
		agg.SpreadPlus25, agg.SpreadPlus25Lower95, agg.SpreadPlus25Upper95 = spreadProb(spreads, -2.5)
		agg.SpreadPlus35, agg.SpreadPlus35Lower95, agg.SpreadPlus35Upper95 = spreadProb(spreads, -3.5)
		agg.SpreadPlus45, agg.SpreadPlus45Lower95, agg.SpreadPlus45Upper95 = spreadProb(spreads, -4.5)
		agg.SpreadPlus55, agg.SpreadPlus55Lower95, agg.SpreadPlus55Upper95 = spreadProb(spreads, -5.5)

		// Moneyline
		mlHomeWin := 0
		for _, s := range spreads {
			if s > 0 {
				mlHomeWin++
			}
		}
		agg.MoneylineHomeWin = float64(mlHomeWin) / float64(len(spreads))
		agg.MlHomeWinLower95, agg.MlHomeWinUpper95 = utils.BinomialCI(agg.MoneylineHomeWin, len(spreads))

		mlAwayWin := 0
		for _, s := range spreads {
			if s < 0 {
				mlAwayWin++
			}
		}
		agg.MoneylineAwayWin = float64(mlAwayWin) / float64(len(spreads))

		// Quantiles, stdev, etc (for all runs, home, away, spread)
		countsFloat := func(vals []int) []float64 {
			out := make([]float64, len(vals))
			for i, v := range vals {
				out[i] = float64(v)
			}
			return out
		}
		agg.StdTotalRuns = utils.Stddev(countsFloat(totalRuns))
		agg.IqrTotalRuns = utils.QuantileWidth(countsFloat(totalRuns), 0.25, 0.75)
		agg.Q80TotalRuns = utils.QuantileWidth(countsFloat(totalRuns), 0.10, 0.90)

		agg.StdHomeScore = utils.Stddev(countsFloat(homeScores))
		agg.IqrHomeScore = utils.QuantileWidth(countsFloat(homeScores), 0.25, 0.75)
		agg.Q80HomeScore = utils.QuantileWidth(countsFloat(homeScores), 0.10, 0.90)
		_, agg.HomeScoreLower95, agg.HomeScoreUpper95 = utils.MeanCI(countsFloat(homeScores))

		agg.StdAwayScore = utils.Stddev(countsFloat(awayScores))
		agg.IqrAwayScore = utils.QuantileWidth(countsFloat(awayScores), 0.25, 0.75)
		agg.Q80AwayScore = utils.QuantileWidth(countsFloat(awayScores), 0.10, 0.90)
		_, agg.AwayScoreLower95, agg.AwayScoreUpper95 = utils.MeanCI(countsFloat(awayScores))

		agg.StdSpread = utils.Stddev(countsFloat(spreads))
		agg.IqrSpread = utils.QuantileWidth(countsFloat(spreads), 0.25, 0.75)
		agg.Q80Spread = utils.QuantileWidth(countsFloat(spreads), 0.10, 0.90)
		_, agg.SpreadLower95, agg.SpreadUpper95 = utils.MeanCI(countsFloat(spreads))

		// --- Averages ---
		agg.AvgTotalRuns = utils.Mean(countsFloat(totalRuns))
		agg.AvgHomeScore = utils.Mean(countsFloat(homeScores))
		agg.AvgAwayScore = utils.Mean(countsFloat(awayScores))

		// --- PLAYER-LEVEL AGGREGATION ---

		batterOutCounts := map[int64][]int{}
		batterKCounts := map[int64][]int{}
		batterBBCounts := map[int64][]int{}
		batterHitCounts := map[int64][]int{}
		batterSingleCounts := map[int64][]int{}
		batterDoubleCounts := map[int64][]int{}
		batterTripleCounts := map[int64][]int{}
		batterHomerunCounts := map[int64][]int{}
		batterRunCounts := map[int64][]int{}
		batterRBICounts := map[int64][]int{}

		// Pitcher per-sim distributions (derived from BatterEvents)
		pitcherKCounts := map[int64][]int{}
		pitcherOutCounts := map[int64][]int{}
		pitcherBBCounts := map[int64][]int{}
		pitcherHitCounts := map[int64][]int{}
		pitcherSwStrCounts := map[int64][]int{}
		pitcherPitchCounts := map[int64][]int{}

		// --- Build stat counts per player for each sim ---
		for _, sim := range outcomes {
			// ---- Batter aggregation ----
			batterCountMap := map[int64]map[string]int{}
			rbiMap := map[int64]int{}
			runMap := map[int64]int{}

			for _, e := range sim.BatterEvents {
				if batterCountMap[e.BatterID] == nil {
					batterCountMap[e.BatterID] = map[string]int{}
				}
				batterCountMap[e.BatterID][e.EventType]++
				switch e.EventType {
				case "single", "double", "triple", "home_run":
					batterCountMap[e.BatterID]["hits"]++
				}
				rbiMap[e.BatterID] += e.RBI
			}
			for _, r := range sim.RunEvents {
				runMap[r.RunnerID]++
			}
			for batterID, m := range batterCountMap {
				batterOutCounts[batterID] = append(batterOutCounts[batterID], m["out"])
				batterKCounts[batterID] = append(batterKCounts[batterID], m["strikeout"])
				batterBBCounts[batterID] = append(batterBBCounts[batterID], m["walk"])
				batterHitCounts[batterID] = append(batterHitCounts[batterID], m["hits"])
				batterSingleCounts[batterID] = append(batterSingleCounts[batterID], m["single"])
				batterDoubleCounts[batterID] = append(batterDoubleCounts[batterID], m["double"])
				batterTripleCounts[batterID] = append(batterTripleCounts[batterID], m["triple"])
				batterHomerunCounts[batterID] = append(batterHomerunCounts[batterID], m["home_run"])
			}
			for id, v := range rbiMap {
				batterRBICounts[id] = append(batterRBICounts[id], v)
			}
			for id, v := range runMap {
				batterRunCounts[id] = append(batterRunCounts[id], v)
			}

			// ---- Pitcher aggregation (derive from BatterEvents) ----
			kPerSim := map[int64]int{}
			outsPerSim := map[int64]int{}
			walksPerSim := map[int64]int{}
			hitsPerSim := map[int64]int{}
			appeared := map[int64]bool{}
			swstrPerSim := map[int64]int{}
			pitchesPerSim := map[int64]int{}

			for _, e := range sim.BatterEvents {
				pid := e.PitcherID
				appeared[pid] = true

				switch e.EventType {
				case "strikeout":
					kPerSim[pid]++
					outsPerSim[pid]++ // strikeout counts as an out
				case "out":
					outsPerSim[pid]++
				case "walk":
					walksPerSim[pid]++
				case "single", "double", "triple", "home_run":
					hitsPerSim[pid]++
				}

				swstrPerSim[e.PitcherID] += int(e.SwStrCount)
				pitchesPerSim[e.PitcherID] += int(e.PitchCount)
			}

			// Append one value per appeared pitcher (zeros if nothing recorded this sim)
			for pid := range appeared {
				pitcherKCounts[pid] = append(pitcherKCounts[pid], kPerSim[pid])
				pitcherOutCounts[pid] = append(pitcherOutCounts[pid], outsPerSim[pid])
				pitcherBBCounts[pid] = append(pitcherBBCounts[pid], walksPerSim[pid])
				pitcherHitCounts[pid] = append(pitcherHitCounts[pid], hitsPerSim[pid])
				pitcherSwStrCounts[pid] = append(pitcherSwStrCounts[pid], swstrPerSim[pid])
				pitcherPitchCounts[pid] = append(pitcherPitchCounts[pid], pitchesPerSim[pid])
			}
		}

		// --- Upsert batter props ---
		for batterID, hitCounts := range batterHitCounts {
			// Aggregate ALL batter stats
			outCounts := batterOutCounts[batterID]
			kCounts := batterKCounts[batterID]
			bbCounts := batterBBCounts[batterID]
			singleCounts := batterSingleCounts[batterID]
			doubleCounts := batterDoubleCounts[batterID]
			tripleCounts := batterTripleCounts[batterID]
			hrCounts := batterHomerunCounts[batterID]
			runCounts := batterRunCounts[batterID]
			rbiCounts := batterRBICounts[batterID]
			// totalHits := len(singleCounts) + len(doubleCounts) + len(tripleCounts) + len(hrCounts)

			prob05Hits, prob15Hits, avgHits, _, iqrHits, q80Hits, lower05Hits, upper05Hits, lower15Hits, upper15Hits := aggregateEventCounts(hitCounts)
			prob05Singles, prob15Singles, avgSingles, _, iqrSingles, q80Singles, lower05Singles, upper05Singles, lower15Singles, upper15Singles := aggregateEventCounts(singleCounts)
			prob05Doubles, prob15Doubles, avgDoubles, _, iqrDoubles, q80Doubles, lower05Doubles, upper05Doubles, lower15Doubles, upper15Doubles := aggregateEventCounts(doubleCounts)
			prob05Triples, prob15Triples, avgTriples, _, iqrTriples, q80Triples, lower05Triples, upper05Triples, lower15Triples, upper15Triples := aggregateEventCounts(tripleCounts)
			prob05HR, prob15HR, avgHR, _, iqrHR, q80HR, lower05HR, upper05HR, lower15HR, upper15HR := aggregateEventCounts(hrCounts)
			prob05Runs, prob15Runs, _, totalRuns, iqrRuns, q80Runs, lower05Runs, upper05Runs, lower15Runs, upper15Runs := aggregateEventCounts(runCounts)
			prob05RBI, prob15RBI, avgRBI, _, iqrRBI, q80RBI, lower05RBI, upper05RBI, lower15RBI, upper15RBI := aggregateEventCounts(rbiCounts)

			totalHits := 0
			for _, h := range hitCounts {
				totalHits += h
			}
			totalSingles := utils.Sum(singleCounts)
			totalDoubles := utils.Sum(doubleCounts)
			totalTriples := utils.Sum(tripleCounts)
			totalHR := utils.Sum(hrCounts)
			totalStrikeouts := utils.Sum(kCounts)
			totalWalks := utils.Sum(bbCounts)

			avgRuns := totalRuns / float64(body.NSims)
			// avgRBI := totalRBI / float64(body.NSims)

			totalAtBats := utils.Sum(hitCounts) + utils.Sum(kCounts) + utils.Sum(outCounts)
			totalPlateAppearances := totalAtBats + totalWalks

			battingAvg := 0.0
			sluggingPct := 0.0
			onBasePct := 0.0
			kPct := 0.0
			bbPct := 0.0
			if totalAtBats > 0 {
				battingAvg = float64(totalHits) / float64(totalAtBats)

				onBasePct = (float64(totalHits) + float64(utils.Sum(bbCounts))) / float64(totalAtBats)

				kPct = float64(totalStrikeouts) / float64(totalPlateAppearances)
				bbPct = float64(totalWalks) / float64(totalPlateAppearances)

				totalBases :=
					1*totalSingles +
						2*totalDoubles +
						3*totalTriples +
						4*totalHR
				sluggingPct = float64(totalBases) / float64(totalAtBats)
			}

			meta := batterMeta[batterID]
			bp := models.BatterProps{
				BatterID:       batterID,
				GamePk:         int64(data.GamePk),
				NumSimulations: len(hitCounts),
				// TODO: Fill these from your player lookup if you have it:
				BatterName:   meta.Name,  // lookupBatterName(batterID)
				BattingOrder: meta.Order, // lookupBattingOrder(batterID)
				Team:         meta.Team,  // lookupTeam(batterID)
				GameDate:     gameDate,   // set from sim/game context

				// Hits
				ProbOver05Hits: prob05Hits, ProbOver15Hits: prob15Hits,
				Over05HitsLower95: lower05Hits, Over05HitsUpper95: upper05Hits,
				Over15HitsLower95: lower15Hits, Over15HitsUpper95: upper15Hits,
				AvgHits: avgHits, IqrHits: iqrHits, Q80Hits: q80Hits,

				// Singles
				ProbOver05Singles: prob05Singles, ProbOver15Singles: prob15Singles,
				Over05SinglesLower95: lower05Singles, Over05SinglesUpper95: upper05Singles,
				Over15SinglesLower95: lower15Singles, Over15SinglesUpper95: upper15Singles,
				AvgSingles: avgSingles, IqrSingles: iqrSingles, Q80Singles: q80Singles,

				// Doubles
				ProbOver05Doubles: prob05Doubles, ProbOver15Doubles: prob15Doubles,
				Over05DoublesLower95: lower05Doubles, Over05DoublesUpper95: upper05Doubles,
				Over15DoublesLower95: lower15Doubles, Over15DoublesUpper95: upper15Doubles,
				AvgDoubles: avgDoubles, IqrDoubles: iqrDoubles, Q80Doubles: q80Doubles,

				// Triples
				ProbOver05Triples: prob05Triples, ProbOver15Triples: prob15Triples,
				Over05TriplesLower95: lower05Triples, Over05TriplesUpper95: upper05Triples,
				Over15TriplesLower95: lower15Triples, Over15TriplesUpper95: upper15Triples,
				AvgTriples: avgTriples, IqrTriples: iqrTriples, Q80Triples: q80Triples,

				// Homeruns
				ProbOver05Homeruns: prob05HR, ProbOver15Homeruns: prob15HR,
				Over05HomerunsLower95: lower05HR, Over05HomerunsUpper95: upper05HR,
				Over15HomerunsLower95: lower15HR, Over15HomerunsUpper95: upper15HR,
				AvgHomeruns: avgHR, IqrHomeruns: iqrHR, Q80Homeruns: q80HR,

				ProbOver05Runs: prob05Runs, ProbOver15Runs: prob15Runs,
				Over05RunsLower95: lower05Runs, Over05RunsUpper95: upper05Runs,
				Over15RunsLower95: lower15Runs, Over15RunsUpper95: upper15Runs,
				AvgRuns: avgRuns, IqrRuns: iqrRuns, Q80Runs: q80Runs,

				ProbOver05RBI: prob05RBI, ProbOver15RBI: prob15RBI,
				Over05RBILower95: lower05RBI, Over05RBIUpper95: upper05RBI,
				Over15RBILower95: lower15RBI, Over15RBIUpper95: upper15RBI,
				AvgRBI: avgRBI, IqrRBI: iqrRBI, Q80RBI: q80RBI,

				BattingAvg:   battingAvg,
				SluggingPct:  sluggingPct,
				OnBasePct:    onBasePct,
				StrikeoutPct: kPct,
				WalkPct:      bbPct,

				TotalHits:             totalHits,
				TotalSingles:          totalSingles,
				TotalDoubles:          totalDoubles,
				TotalTriples:          totalTriples,
				TotalHomeruns:         totalHR,
				TotalPlateAppearances: totalPlateAppearances,
				TotalAtBats:           totalAtBats,
				TotalStrikeouts:       totalStrikeouts,
				TotalWalks:            totalWalks,
			}
			_ = poster.InsertBatterProps(s.db, bp)
		}

		// --- Upsert pitcher props ---
		for pitcherID, kSlice := range pitcherKCounts {
			// K thresholds (>=3,4,...,13)
			prob25K, l25K, u25K := getKProp(kSlice, 3)
			prob35K, l35K, u35K := getKProp(kSlice, 4)
			prob45K, l45K, u45K := getKProp(kSlice, 5)
			prob55K, l55K, u55K := getKProp(kSlice, 6)
			prob65K, l65K, u65K := getKProp(kSlice, 7)
			prob75K, l75K, u75K := getKProp(kSlice, 8)
			prob85K, l85K, u85K := getKProp(kSlice, 9)
			prob95K, l95K, u95K := getKProp(kSlice, 10)
			prob105K, l105K, u105K := getKProp(kSlice, 11)
			prob115K, l115K, u115K := getKProp(kSlice, 12)
			prob125K, l125K, u125K := getKProp(kSlice, 13)

			// Dispersion stats
			f := make([]float64, len(kSlice))
			for i, k := range kSlice {
				f[i] = float64(k)
			}
			// avgK := utils.Mean(f)
			// fmt.Println("avgK:", kSlice)
			avgK := 0.0
			iqrK := utils.QuantileWidth(f, 0.25, 0.75)
			q80K := utils.QuantileWidth(f, 0.10, 0.90)

			// Totals for rate stats
			outSlice := pitcherOutCounts[pitcherID]
			bbSlice := pitcherBBCounts[pitcherID]
			hitSlice := pitcherHitCounts[pitcherID]
			swstrSlice := pitcherSwStrCounts[pitcherID]
			pitchCountSlice := pitcherPitchCounts[pitcherID]

			totalStrikeouts := utils.Sum(kSlice)
			avgK = float64(totalStrikeouts) / float64(body.NSims)
			totalWalks := utils.Sum(bbSlice)
			totalSwStr := utils.Sum(swstrSlice)
			totalPitches := utils.Sum(pitchCountSlice)

			totalAB := utils.Sum(hitSlice) + totalStrikeouts + utils.Sum(outSlice)
			totalPA := totalAB + totalWalks

			totalOuts := int(totalStrikeouts + utils.Sum(outSlice))

			pitcherKPct := 0.0
			pitcherBBPct := 0.0
			pitcherSwStrPct := 0.0
			inningsPitched := utils.ConvertOutsToInnings(totalOuts)
			if totalPA > 0 {
				pitcherKPct = float64(totalStrikeouts) / float64(totalPA)
				pitcherBBPct = float64(totalWalks) / float64(totalPA)
				pitcherSwStrPct = float64(totalSwStr) / float64(totalPitches)
			}

			meta := pitcherMeta[pitcherID]
			pp := models.PitcherProps{
				PitcherID:      pitcherID,
				GamePk:         int64(data.GamePk),
				NumSimulations: len(kSlice),
				PitcherName:    meta.Name,
				GameDate:       gameDate,
				Team:           meta.Team,

				ProbOver25K: prob25K, ProbOver35K: prob35K, ProbOver45K: prob45K,
				ProbOver55K: prob55K, ProbOver65K: prob65K, ProbOver75K: prob75K,
				ProbOver85K: prob85K, ProbOver95K: prob95K, ProbOver105K: prob105K,
				ProbOver115K: prob115K, ProbOver125K: prob125K,
				Over25KLower95: l25K, Over25KUpper95: u25K,
				Over35KLower95: l35K, Over35KUpper95: u35K,
				Over45KLower95: l45K, Over45KUpper95: u45K,
				Over55KLower95: l55K, Over55KUpper95: u55K,
				Over65KLower95: l65K, Over65KUpper95: u65K,
				Over75KLower95: l75K, Over75KUpper95: u75K,
				Over85KLower95: l85K, Over85KUpper95: u85K,
				Over95KLower95: l95K, Over95KUpper95: u95K,
				Over105KLower95: l105K, Over105KUpper95: u105K,
				Over115KLower95: l115K, Over115KUpper95: u115K,
				Over125KLower95: l125K, Over125KUpper95: u125K,
				AvgStrikeouts: avgK, IqrStrikeouts: iqrK, Q80Strikeouts: q80K,

				StrikeoutPct:          pitcherKPct,
				WalkPct:               pitcherBBPct,
				SwingingStrPct:        pitcherSwStrPct,
				TotalStrikeouts:       totalStrikeouts,
				TotalWalks:            totalWalks,
				TotalPlateAppearances: totalPA,
				TotalSwingingStrikes:  totalSwStr,
				TotalPitches:          totalPitches,
				IP:                    inningsPitched,
			}
			if err := poster.InsertPitcherProps(s.db, pp); err != nil {
				log.Printf("InsertPitcherProps error: %v", err)
			}
		}

		// --- UPSERT TO DB ---
		if err := poster.InsertGameResultAggCore(s.db, agg); err != nil {
			log.Printf("Insert agg result failed: %v", err)
		}

		_, err = s.db.Exec(`
                       UPDATE simulation_jobs
                       SET status = 'completed', result = $1, current_simulation = total_simulations, updated_at = NOW()
                       WHERE id = $2`, "Simulations complete", jid)
		if err != nil {
			log.Printf("Failed to update job %s to completed: %v", jid, err)
		}
	}(jobID, body.UserID, gameData, body.NSims, body.OutcomeMode)

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
	var currentSim, totalSim int

	err := s.db.QueryRow(`
               SELECT status, result, current_simulation, total_simulations FROM simulation_jobs
               WHERE id = $1 AND user_id = $2`, body.JobID, body.UserID).Scan(&status, &result, &currentSim, &totalSim)
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
	json.NewEncoder(w).Encode(map[string]interface{}{
		"jobId":             body.JobID,
		"status":            status,
		"result":            res,
		"currentSimulation": currentSim,
		"totalSimulations":  totalSim,
	})
}

func pacificDayRangeUTC(isoDate string) (startUTC, endUTC time.Time, err error) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	// Parse YYYY-MM-DD as a *local* Pacific midnight
	t, err := time.ParseInLocation("2006-01-02", isoDate, loc)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	start := t
	end := t.Add(24 * time.Hour)

	return start.UTC(), end.UTC(), nil
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
		var currentSim, totalSim int

		err := s.db.QueryRow(`
                       SELECT status, result, current_simulation, total_simulations FROM simulation_jobs
                       WHERE id = $1 AND user_id = $2`, jobID, userID).Scan(&status, &result, &currentSim, &totalSim)
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
		json.NewEncoder(w).Encode(map[string]interface{}{
			"jobId":             jobID,
			"status":            status,
			"result":            res,
			"currentSimulation": currentSim,
			"totalSimulations":  totalSim,
		})
		return
	}

	// Case 2: Return all jobs for user
	rows, err := s.db.Query(`
               SELECT id, status, result, current_simulation, total_simulations FROM simulation_jobs
               WHERE user_id = $1 ORDER BY updated_at DESC`, userID)
	if err != nil {
		log.Printf("Error querying jobs for user %s: %v", userID, err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var jobs []map[string]interface{}

	for rows.Next() {
		var id, status string
		var result sql.NullString
		var currentSim, totalSim int

		if err := rows.Scan(&id, &status, &result, &currentSim, &totalSim); err != nil {
			log.Printf("Error scanning job row: %v", err)
			continue
		}

		res := ""
		if result.Valid {
			res = result.String
		}

		jobs = append(jobs, map[string]interface{}{
			"jobId":             id,
			"status":            status,
			"result":            res,
			"currentSimulation": currentSim,
			"totalSimulations":  totalSim,
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
		startUTC, endUTC, err := pacificDayRangeUTC(gameDateStr)
		if err != nil {
			http.Error(w, "invalid gamedate", http.StatusBadRequest)
			return
		}
		where = append(where, fmt.Sprintf("gamedate >= $%d AND gamedate < $%d", argID, argID+1))
		args = append(args, startUTC, endUTC)
		argID += 2
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
