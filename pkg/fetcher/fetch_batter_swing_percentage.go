package fetcher

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/logananthony/go-baseball/pkg/models"
)


// FetchPitcherFrequencies queries and prints pitch data for a pitcher ID
func FetchBatterSwingPercentage(db *sql.DB, batterId, gameYear int) ([]models.BatterSwingPercentage, error) {
	query := `
		SELECT game_year, batter, stand, p_throws, zone, pitch_type, total_pitches, total_swings, swing_percentage
		FROM batter_swing_percentage
		WHERE batter = $1 AND game_year = $2;
	`

	rows, err := db.Query(query, batterId, gameYear)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var results []models.BatterSwingPercentage

	for rows.Next() {
		var rec models.BatterSwingPercentage
		err := rows.Scan(
			&rec.GameYear,
			&rec.Batter,
			&rec.Stand,
			&rec.PThrows,
			&rec.Zone,
			&rec.PitchType,
			&rec.TotalPitches,
			&rec.TotalSwings,
			&rec.SwingPercentage,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		results = append(results, rec)
	}

	return results, nil
}



