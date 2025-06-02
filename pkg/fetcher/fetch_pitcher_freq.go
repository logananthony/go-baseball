package fetcher

import (
	"database/sql"
	//"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/logananthony/go-baseball/pkg/models"
)

// FetchPitcherFrequencies queries and prints pitch data for a pitcher ID
func FetchPitcherFrequencies(db *sql.DB, pitcherID int, stand string) []models.PitcherCountPitchFreq {
	query := `SELECT stand, pitch_type, balls, strikes, count, frequency FROM pitcher_count_pitch_freq WHERE pitcher = $1 AND stand = $2`

	rows, err := db.Query(query, pitcherID, stand)
	if err != nil {
		log.Fatal("Query error:", err)
	}
	defer rows.Close()

	var results []models.PitcherCountPitchFreq

	for rows.Next() {
		var freq models.PitcherCountPitchFreq
		freq.PITCHER = pitcherID
		freq.STAND = stand
		err := rows.Scan(
			&freq.STAND,
			&freq.PITCH_TYPE,
			&freq.BALLS,
			&freq.STRIKES,
			&freq.COUNT,
			&freq.FREQUENCY,
		)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, freq)
	}

	return results
}
