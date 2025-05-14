package fetcher

import (
	"database/sql"
	"log"

	"github.com/logananthony/go-baseball/pkg/models"
)

// FetchPitchingSubstitutionProbs retrieves all rows from the pitching_substitution_probs table
func FetchPitchingSubstitutionProbs(db *sql.DB) ([]models.PitchingSubstitutionProb, error) {
	query := `
		SELECT
			inning,
			runs_scored_game,
			runs_scored_inning,
			total_appearances,
			total_pulled,
			pull_probability
		FROM pitching_substitution_probs
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.PitchingSubstitutionProb

	for rows.Next() {
		var record models.PitchingSubstitutionProb
		if err := rows.Scan(
			&record.Inning,
			&record.RunsScoredGame,
			&record.RunsScoredInning,
			&record.TotalAppearances,
			&record.TotalPulled,
			&record.PullProbability,
		); err != nil {
			log.Println("Scan error:", err)
			continue
		}
		results = append(results, record)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

