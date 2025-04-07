package fetcher

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/logananthony/go-baseball/pkg/models"
)

// FetchPlayerInfo retrieves a player's info by ID, full name, and/or season (any can be nil/zero).
func FetchPlayerInfo(db *sql.DB, id *int, fullName *string, season *int) ([]models.MLBPlayerInfo, error) {
	var (
		clauses []string
		args    []interface{}
	)

	if id != nil {
		clauses = append(clauses, "id = $"+fmt.Sprint(len(args)+1))
		args = append(args, *id)
	}
	if fullName != nil {
		clauses = append(clauses, "LOWER(fullName) = LOWER($"+fmt.Sprint(len(args)+1)+")")
		args = append(args, *fullName)
	}
	if season != nil {
		clauses = append(clauses, "season = $"+fmt.Sprint(len(args)+1))
		args = append(args, *season)
	}

	if len(clauses) == 0 {
		return nil, fmt.Errorf("must supply at least one filter: id, fullName, or season")
	}

	query := `
SELECT 
	id, fullName, firstName, lastName, primaryNumber, birthDate, currentAge,
	birthCity, birthStateProvince, birthCountry, height, weight, active,
	batSide, pitchHand, mlbDebutDate, strikeZoneTop, strikeZoneBottom,
	teamId, teamName, position, season
FROM mlb_player_info
WHERE ` + strings.Join(clauses, " AND ")

	rows, err := db.Query(query, args...)
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

