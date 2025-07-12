// In fetcher/fetcher.go

package fetcher

import (
	"database/sql"

	"github.com/logananthony/go-baseball/pkg/models"
)

func FetchBullpenRoleProbs(db *sql.DB, homeTeam string, awayTeam string, season int) ([]models.BullpenRoleProb, error) {
	rows, err := db.Query(`
		SELECT inning, run_diff, runners_on, role, probability
		FROM bullpen_role_probs
		WHERE (team_abbr = $1 OR team_abbr = $2) AND season = $3
	`, homeTeam, awayTeam, season)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.BullpenRoleProb
	for rows.Next() {
		var rp models.BullpenRoleProb
		if err := rows.Scan(&rp.Inning, &rp.RunDiff, &rp.RunnersOn, &rp.Role, &rp.Prob); err != nil {
			return nil, err
		}
		result = append(result, rp)
	}
	return result, nil
}
