package fetcher

import (
	"database/sql"
	"log"

	"github.com/logananthony/go-baseball/pkg/models"
)

// FetchAllPitcherExitProbss returns all rows in pitcher_exit_distribution.
func FetchPitchingSubstitutionProbs(db *sql.DB) ([]models.PitchingSubstitutionProb, error) {
	const q = `
		SELECT
			pitcher_role,
			pitch_bin,
			num_pitchers,
			pct_of_role,
			cumulative_pct_of_role
		FROM pitcher_exit_probs
		ORDER BY pitcher_role, pitch_bin;
	`
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.PitchingSubstitutionProb
	for rows.Next() {
		var ped models.PitchingSubstitutionProb
		err := rows.Scan(
			&ped.PitcherRole,
			&ped.PitchBin,
			&ped.NumPitchers,
			&ped.PctOfRole,
			&ped.CumulativePctOfRole,
		)
		if err != nil {
			log.Printf("fetcher.FetchAllPitcherExitProbss scan error: %v", err)
			continue
		}
		results = append(results, ped)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
