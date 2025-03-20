package fetcher

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/logananthony/go-baseball/pkg/models"
)


// FetchPitcherFrequencies queries and prints pitch data for a pitcher ID
func FetchBatterContactPercentage(db *sql.DB, batterId, gameYear int) ([]models.BatterContactPercentage, error) {
	query := `
		SELECT game_year, batter, stand, p_throws, zone, pitch_type, swinging_strike, foul, ball_in_play, total_swings, pct_swinging_strike, pct_foul, pct_ball_in_play
		FROM batter_contact_percentage
		WHERE batter = $1 AND game_year = $2;
	`

	rows, err := db.Query(query, batterId, gameYear)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var results []models.BatterContactPercentage

	for rows.Next() {
		var rec models.BatterContactPercentage
		err := rows.Scan(
			&rec.GameYear,
			&rec.Batter,
			&rec.Stand,
			&rec.PThrows,
			&rec.Zone,
			&rec.PitchType,
			&rec.SwingingStrike,
			&rec.Foul,
			&rec.BallInPlay,
      &rec.TotalSwings, 
      &rec.PctSwingingStrike, 
      &rec.PctFoul, 
      &rec.PctBallInPlay,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		results = append(results, rec)
	}

	return results, nil
}




