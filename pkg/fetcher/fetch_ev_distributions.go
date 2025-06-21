package fetcher

import (
	"database/sql"
	"log"

	"github.com/logananthony/go-baseball/pkg/models"
)

func FetchEVDistributions(
	db *sql.DB,
	gameYear int,
	batter int,
) []models.EVDistribution {

	query := `
SELECT
	game_year,
	batter,
	stand,
	p_throws,
	outcome,
	pitch_type,
	zone,
	velocity_bucket,
	skew,
	mean,
	std,
	n,
	level
FROM ev_distributions
WHERE
	game_year = $1 AND
	batter = $2
`

	rows, err := db.Query(query, gameYear, batter)
	if err != nil {
		log.Fatal("Query error:", err)
	}
	defer rows.Close()

	var results []models.EVDistribution
	for rows.Next() {
		var ev models.EVDistribution
		err := rows.Scan(
			&ev.GameYear,
			&ev.Batter,
			&ev.Stand,
			&ev.PThrows,
			&ev.Outcome,
			&ev.PitchType,
			&ev.Zone,
			&ev.VelocityBucket,
			&ev.Skew,
			&ev.Mean,
			&ev.Std,
			&ev.N,
			&ev.Level,
		)
		if err != nil {
			log.Fatal("Scan error:", err)
		}

		// Debug logging for each row
		log.Printf("Fetched EV row: mean.Valid=%v, std.Valid=%v, skew.Valid=%v | mean=%.2f, std=%.2f",
			ev.Mean.Valid, ev.Std.Valid, ev.Skew.Valid,
			ev.Mean.Float64, ev.Std.Float64,
		)

		// Optional: skip invalid rows
		if !ev.Mean.Valid || !ev.Std.Valid {
			log.Printf("Skipping EV row due to invalid mean/std for batter %d", batter)
			continue
		}

		results = append(results, ev)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Row iteration error:", err)
	}

	if len(results) == 0 {
		log.Printf("No EV distributions found for batter %d in year %d", batter, gameYear)
	}

	return results
}
