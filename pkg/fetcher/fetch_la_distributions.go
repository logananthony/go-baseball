package fetcher

import (
	"database/sql"
	"log"

	"github.com/logananthony/go-baseball/pkg/models"
)

func FetchLADistributions(
	db *sql.DB,
	gameYear int,
	batter int,
) []models.LADistribution {

	query := `
SELECT
	game_year,
	batter,
	stand,
	p_throws,
	outcome,
	zone,
	ev_bucket,
	skew,
	mean,
	std,
	n,
	level
FROM la_distributions
WHERE
	game_year = $1 AND
	batter = $2 
	`

	rows, err := db.Query(
		query,
		gameYear,
		batter,
	)
	if err != nil {
		log.Fatal("Query error:", err)
	}
	defer rows.Close()

	var results []models.LADistribution
	for rows.Next() {
		var ev models.LADistribution
		err := rows.Scan(
			&ev.GameYear,
			&ev.Batter,
			&ev.Stand,
			&ev.PThrows,
			&ev.Outcome,
			&ev.Zone,
			&ev.EVBucket,
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
