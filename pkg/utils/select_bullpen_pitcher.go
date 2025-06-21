package utils

import (
	"database/sql"
	"log"
	"math/rand"
)

func SelectBullpenPitcherLineup(
	db *sql.DB,
	pitcherLineup [][]int,
	inning int,
	runDiff int,
	runnersOn int,
	used map[int]bool,
) []int {
	query := `
		SELECT role, probability
		FROM bullpen_role_probs
		WHERE inning = $1 AND run_diff = $2 AND runners_on = $3
	`
	rows, err := db.Query(query, inning, runDiff, runnersOn)
	if err != nil {
		log.Printf("❌ Failed to query bullpen_role_probs: %v", err)
		return nil
	}
	defer rows.Close()

	var roles []int
	var probs []float64

	for rows.Next() {
		var role int
		var prob float64
		if err := rows.Scan(&role, &prob); err != nil {
			log.Printf("❌ Error scanning row: %v", err)
			continue
		}
		roles = append(roles, role)
		probs = append(probs, prob)
	}

	if len(roles) == 0 {
		log.Printf("⚠️ No bullpen data found for inning=%d run_diff=%d runners_on=%d", inning, runDiff, runnersOn)
		return nil
	}

	var eligible [][]int
	var weights []float64

	for i, roleIdx := range roles {
		if roleIdx >= len(pitcherLineup) {
			continue
		}
		pid := pitcherLineup[roleIdx][0]
		if !used[pid] {
			eligible = append(eligible, pitcherLineup[roleIdx])
			weights = append(weights, probs[i])
		}
	}

	if len(eligible) == 0 {
		return nil
	}

	// Normalize weights
	total := 0.0
	for _, w := range weights {
		total += w
	}
	for i := range weights {
		weights[i] /= total
	}

	// Weighted random choice
	r := rand.Float64()
	cumulative := 0.0
	for i, prob := range weights {
		cumulative += prob
		if r <= cumulative {
			return eligible[i]
		}
	}

	return eligible[len(eligible)-1]
}
