package fetcher

import (
	"database/sql"
	"log"

	"github.com/logananthony/go-baseball/pkg/models"
)

// FetchAllBullpenRoleProbs fetches all bullpen role probabilities.
func FetchAllBullpenRoleProbs(db *sql.DB) ([]models.BullpenRoleProb, error) {
	const q = `
		SELECT
			inning,
			run_diff,
			runners_on,
			role,
			probability
		FROM bullpen_role_probs;
	`
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.BullpenRoleProb
	for rows.Next() {
		var rp models.BullpenRoleProb
		err := rows.Scan(
			&rp.Inning,
			&rp.RunsScoredGame,
			&rp.RunsScoredInning,
			&rp.Role,
			&rp.PullProbability,
		)
		if err != nil {
			log.Printf("fetcher.FetchAllBullpenRoleProbs scan error: %v", err)
			continue
		}
		results = append(results, rp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
