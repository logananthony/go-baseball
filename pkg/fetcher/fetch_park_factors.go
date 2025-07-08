package fetcher

import (
	"database/sql"
	"log"

	"github.com/logananthony/go-baseball/pkg/models"
)

// FetchParkFactors returns park-factor data for the given team abbreviation.
//
// The first return value is the struct itself (non-pointer).
// The second return value is true if a row was found, false otherwise.
func FetchParkFactors(db *sql.DB, teamAbbr string) (models.ParkFactors, bool) {
	// ─── Alias map (tweak if needed) ───────────────────────────────────────
	abbrMap := map[string]string{
		"AZ":  "ARI",
		"ATH": "OAK",
		// 	"KC":  "KCR",
		// 	"SF":  "SFG",
		// 	"TB":  "TBR",
		// 	"WSH": "WSN",
	}
	if val, ok := abbrMap[teamAbbr]; ok {
		teamAbbr = val
	}

	const q = `
SELECT
	"Team",
	"Venue",
	"H",
	"1B",
	"2B",
	"3B",
	"HR",
	"OBP"
FROM park_factors
WHERE UPPER("Team") = UPPER($1)
LIMIT 1;`

	row := db.QueryRow(q, teamAbbr)

	var pf models.ParkFactors
	if err := row.Scan(
		&pf.Team,
		&pf.Venue,
		&pf.H,
		&pf.OneB,
		&pf.TwoB,
		&pf.ThreeB,
		&pf.HR,
		&pf.OBP,
	); err != nil {
		if err == sql.ErrNoRows {
			return models.ParkFactors{}, false // not found
		}
		log.Fatalf("fetcher.FetchParkFactors: %v", err)
	}

	return pf, true
}
