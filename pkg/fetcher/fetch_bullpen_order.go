package fetcher

import (
	"database/sql"
	"log"

	"github.com/logananthony/go-baseball/pkg/models"
)

func FetchBullpenOrder(db *sql.DB, teamAbbr string, season int) *models.BullpenOrder {

	var teamAbbrMap = map[string]string{
		"AZ": "ARI",
		// "SF": "SFG",
		// "TB": "TBR",
		// "KC": "KCR",
		// "WSH": "WSN",
		// Add any others used in your data sources
	}

	if val, ok := teamAbbrMap[teamAbbr]; ok {
		teamAbbr = val
	}

	query := `
SELECT 
	season,
	abbreviation,
	roster_team_id,
	player_id_1,
	player_id_2,
	player_id_3,
	player_id_4,
	player_id_5,
	player_id_6,
	player_id_7,
	player_id_8
FROM bullpen_order
WHERE LOWER(abbreviation) = LOWER($1) AND season = $2
LIMIT 1;
`

	row := db.QueryRow(query, teamAbbr, season)

	var bp models.BullpenOrder
	var ignoredSeason int // we don't need to keep this

	err := row.Scan(
		&ignoredSeason, // ignore this field
		&bp.TeamAbbreviation,
		&bp.RosterTeamID,
		&bp.PlayerID1,
		&bp.PlayerID2,
		&bp.PlayerID3,
		&bp.PlayerID4,
		&bp.PlayerID5,
		&bp.PlayerID6,
		&bp.PlayerID7,
		&bp.PlayerID8,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatalf("Failed to scan bullpen order: %v", err)
	}

	return &bp
}
