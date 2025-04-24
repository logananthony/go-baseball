package fetcher

import (
	"database/sql"
	"log"
	"github.com/logananthony/go-baseball/pkg/models"
)

func FetchSprayDistributions(
	db *sql.DB,
	gameYear int,
	batter int,
	stand string,
	pThrows string,
	outcome sql.NullString,
	zone sql.NullInt32,
	evBucket sql.NullString,
	launchAngleBucket sql.NullString,
) []models.SprayDistribution {

	query := `
SELECT
	game_year,
	batter,
	stand,
	p_throws,
	outcome,
	zone,
	ev_bucket,
	launch_angle_bucket,
	skew,
	mean,
	std,
	n,
	level
FROM spray_distributions
WHERE
	game_year = $1 AND
	batter = $2 AND
	stand = $3 AND
	p_throws = $4 AND
	(outcome IS NULL OR outcome = $5) AND
	(zone IS NULL OR zone = $6) AND
	(ev_bucket IS NULL OR ev_bucket = $7) AND
	(launch_angle_bucket IS NULL OR launch_angle_bucket = $8)
	`

	rows, err := db.Query(
		query,
		gameYear,
		batter,
		stand,
		pThrows,
		outcome,
		zone,
		evBucket,
		launchAngleBucket,
	)
	if err != nil {
		log.Fatal("Query error:", err)
	}
	defer rows.Close()

	var results []models.SprayDistribution
	for rows.Next() {
		var sd models.SprayDistribution
		err := rows.Scan(
			&sd.GameYear,
			&sd.Batter,
			&sd.Stand,
			&sd.PThrows,
			&sd.Outcome,
			&sd.Zone,
			&sd.EVBucket,
			&sd.LaunchAngleBucket,
			&sd.Skew,
			&sd.Mean,
			&sd.Std,
			&sd.N,
			&sd.Level,
		)
		if err != nil {
			log.Fatal("Scan error:", err)
		}
		results = append(results, sd)
	}

	return results
}

