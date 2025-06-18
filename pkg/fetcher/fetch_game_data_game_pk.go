package fetcher

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/logananthony/go-baseball/pkg/models"
)

// FetchGameDataByGamePk retrieves lineup data for a specific gamePk
func FetchGameDataByGamePk(db *sql.DB, gamePk int) ([]models.GameDataGamePk, error) {
	query := `
SELECT 
    gamePk,
    season,
    gameDate,
    team,
    teamAbbreviation,
    battingOrder,
    playerId,
    playerName,
    homePitcherId,
    homePitcherName,
    awayPitcherId,
    awayPitcherName,
    homeTeamAbbr,
    awayTeamAbbr
FROM game_data_game_pk
WHERE gamePk = $1;
`

	rows, err := db.Query(query, gamePk)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var results []models.GameDataGamePk

	for rows.Next() {
		var rec models.GameDataGamePk
		var gameDate time.Time // Ensure this matches your struct

		err := rows.Scan(
			&rec.GamePk,
			&rec.Season,
			&gameDate,
			&rec.Team,
			&rec.TeamAbbreviation,
			&rec.BattingOrder,
			&rec.PlayerId,
			&rec.PlayerName,
			&rec.HomePitcherId,
			&rec.HomePitcherName,
			&rec.AwayPitcherId,
			&rec.AwayPitcherName,
			&rec.HomeTeamAbbr,
			&rec.AwayTeamAbbr,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		rec.GameDate = gameDate
		results = append(results, rec)
	}

	return results, nil
}
