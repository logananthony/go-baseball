package fetcher

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/logananthony/go-baseball/pkg/models"
)

// FetchBatterHitType queries batter hit type probabilities from batter_hit_type table
func FetchBatterHitType(db *sql.DB, batterId, gameYear int) ([]models.BatterHitType, error) {
	query := `
		SELECT game_year, batter, stand, p_throws, zone, pitch_type, velocity_bucket,
		       double, home_run, out, single, triple, n, level
		FROM batter_hit_type
		WHERE batter = $1 AND game_year = $2;
	`

	rows, err := db.Query(query, batterId, gameYear)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var results []models.BatterHitType

	for rows.Next() {
		var rec models.BatterHitType
		err := rows.Scan(
			&rec.GameYear,
			&rec.Batter,
			&rec.Stand,
			&rec.PThrows,
			&rec.Zone,
			&rec.PitchType,
			&rec.VelocityBucket,
			&rec.Double,
			&rec.HomeRun,
			&rec.Out,
			&rec.Single,
			&rec.Triple,
			&rec.N,
			&rec.Level,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		results = append(results, rec)
	}

	return results, nil
}


