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
SELECT 
    game_year, 
    batter, 
    COALESCE(stand, 'NA') as stand, 
    COALESCE(p_throws, 'NA') as p_throws, 
    COALESCE(zone, -1::BIGINT) as zone, 
    COALESCE(pitch_type, 'NA') as pitch_type, 
    COALESCE(velocity_bucket, 'NA') as velocity_bucket,
    COALESCE(double, 0.0) as double,
    COALESCE(home_run, 0.0) as home_run,
    COALESCE(out, 0.0) as out,
    COALESCE(single, 0.0) as single,
    COALESCE(triple, 0.0) as triple,
    COALESCE(n, 0) as n,
    COALESCE(level, '-1') as level
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

	//fmt.Printf("Rows fetched: %d\n", len(results))

	return results, nil
}

