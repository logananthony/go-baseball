package fetcher

import (
	"database/sql"
	"fmt"
	//"strings"

	"github.com/logananthony/go-baseball/pkg/models"
)

// FetchPlayerInfo retrieves a player's info by ID, full name, and/or season (any can be nil).
func FetchPlayerInfo(db *sql.DB, id *int, season *int) ([]models.MLBPlayerInfo, error) {
	if id == nil || season == nil {
		return nil, fmt.Errorf("id and season must be provided")
	}

	query := `
SELECT 
	id, fullName, firstName, lastName, primaryNumber, birthDate, currentAge,
	birthCity, birthStateProvince, birthCountry, height, weight, active,
	batSide, pitchHand, mlbDebutDate, strikeZoneTop, strikeZoneBottom,
	teamId, teamName, position, season
FROM mlb_player_info
WHERE id = $1 AND season = $2;
`

//fmt.Printf("Final query: %s\n", query)
//fmt.Printf("Args: %v %v\n", *id, *season)

	rows, err := db.Query(query, *id, *season)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var results []models.MLBPlayerInfo
	for rows.Next() {
		var player models.MLBPlayerInfo
		err := rows.Scan(
			&player.ID,
			&player.FullName,
			&player.FirstName,
			&player.LastName,
			&player.PrimaryNumber,
			&player.BirthDate,
			&player.CurrentAge,
			&player.BirthCity,
			&player.BirthState,
			&player.BirthCountry,
			&player.Height,
			&player.Weight,
			&player.Active,
			&player.BatSide,
			&player.PitchHand,
			&player.MLBDebutDate,
			&player.StrikeZoneTop,
			&player.StrikeZoneBottom,
			&player.TeamID,
			&player.TeamName,
			&player.Position,
			&player.Season,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		results = append(results, player)
	}

	return results, nil
}

