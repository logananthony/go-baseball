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
	stand string,
	pThrows string,
	outcome sql.NullString,
	pitchType sql.NullString,
	zone sql.NullInt32,
	velocityBucket sql.NullString,
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
	batter = $2 AND
	stand = $3 AND
	p_throws = $4 AND
	(outcome is null OR outcome = $5) AND
	(pitch_type is null OR pitch_type = $6) AND
	(zone is null OR zone = $7) AND
	(velocity_bucket is null OR velocity_bucket = $8)
	`

	rows, err := db.Query(
		query,
		gameYear,
		batter,
		stand,
		pThrows,
		outcome,
		pitchType,
		zone,
		velocityBucket,
	)
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
		results = append(results, ev)
	}

	return results
}

// Helper to convert empty string into SQL NULL


