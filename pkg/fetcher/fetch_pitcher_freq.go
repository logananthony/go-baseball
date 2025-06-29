package fetcher

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/logananthony/go-baseball/pkg/models"
)

// // FetchPitcherFrequencies retrieves pitch frequencies for a given pitcher, batter stand, and game year.
// func FetchPitcherFrequencies(db *sql.DB, pitcherID int, stand string, gameYear int) []models.PitcherCountPitchFreq {
// 	query := `
// SELECT
// 	stand,
// 	pitch_type,
// 	balls,
// 	strikes,
// 	count,
// 	frequency
// FROM pitcher_count_pitch_freq
// WHERE pitcher = $1 AND stand = $2 AND game_year = $3;
// `

// 	rows, err := db.Query(query, pitcherID, stand, gameYear)
// 	if err != nil {
// 		log.Fatal("Query error:", err)
// 	}
// 	defer rows.Close()

// 	var results []models.PitcherCountPitchFreq

// 	for rows.Next() {
// 		var freq models.PitcherCountPitchFreq
// 		freq.PITCHER = pitcherID
// 		freq.STAND = stand
// 		freq.GAME_YEAR = gameYear

// 		err := rows.Scan(
// 			&freq.STAND,
// 			&freq.PITCH_TYPE,
// 			&freq.BALLS,
// 			&freq.STRIKES,
// 			&freq.COUNT,
// 			&freq.FREQUENCY,
// 		)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		results = append(results, freq)
// 	}

// 	return results
// }

func FetchPitcherFrequencies(db *sql.DB, pitcherID int, stand string, gameYear int) []models.PitcherCountPitchFreq {
	query := `
SELECT 
	stand, 
	pitch_type, 
	balls, 
	strikes, 
	count, 
	frequency 
FROM pitcher_count_pitch_freq 
WHERE pitcher = $1 AND stand = $2 AND game_year = $3;
`

	rows, err := db.Query(query, pitcherID, stand, gameYear)
	if err != nil {
		log.Fatal("Query error:", err)
	}
	defer rows.Close()

	var results []models.PitcherCountPitchFreq
	found := false

	for rows.Next() {
		found = true
		var freq models.PitcherCountPitchFreq
		freq.PITCHER = pitcherID
		freq.GAME_YEAR = gameYear

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

	if !found {
		results = []models.PitcherCountPitchFreq{}
		pitchTypes := []struct {
			PitchType string
			Frequency float64
		}{
			{"FF", 0.45},
			{"SL", 0.30},
			{"CH", 0.25},
		}

		for balls := 0; balls <= 3; balls++ {
			for strikes := 0; strikes <= 2; strikes++ {
				for _, pt := range pitchTypes {
					results = append(results, models.PitcherCountPitchFreq{
						PITCHER:    pitcherID,
						STAND:      stand,
						PITCH_TYPE: pt.PitchType,
						BALLS:      balls,
						STRIKES:    strikes,
						COUNT:      100,
						FREQUENCY:  pt.Frequency,
						GAME_YEAR:  gameYear,
					})
				}
			}
		}
	}

	return results
}
