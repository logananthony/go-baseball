package poster

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/logananthony/go-baseball/pkg/models"
)

func InsertPitcherProps(db *sql.DB, prop models.PitcherProps) error {
	columns := []string{
		"pitcherid", "gamepk", "pitchername", "gamedate", "num_simulations",
		"prob_over_2_5_k", "prob_over_3_5_k", "prob_over_4_5_k", "prob_over_5_5_k", "prob_over_6_5_k",
		"prob_over_7_5_k", "prob_over_8_5_k", "prob_over_9_5_k", "prob_over_10_5_k", "prob_over_11_5_k", "prob_over_12_5_k",
		"over_2_5_k_lower95", "over_2_5_k_upper95", "over_3_5_k_lower95", "over_3_5_k_upper95",
		"over_4_5_k_lower95", "over_4_5_k_upper95", "over_5_5_k_lower95", "over_5_5_k_upper95",
		"over_6_5_k_lower95", "over_6_5_k_upper95", "over_7_5_k_lower95", "over_7_5_k_upper95",
		"over_8_5_k_lower95", "over_8_5_k_upper95", "over_9_5_k_lower95", "over_9_5_k_upper95",
		"over_10_5_k_lower95", "over_10_5_k_upper95", "over_11_5_k_lower95", "over_11_5_k_upper95",
		"over_12_5_k_lower95", "over_12_5_k_upper95",
		"avg_strikeouts", "iqr_strikeouts", "q80_strikeouts", "team", "strikeout_pct", "walk_pct",
		"swstr_pct", "innings_pitched", "total_strikeouts", "total_walks", "total_plate_appearances",
		"total_swinging_strikes", "total_pitches",
	}

	values := []interface{}{
		prop.PitcherID, prop.GamePk, prop.PitcherName, prop.GameDate, prop.NumSimulations,
		prop.ProbOver25K, prop.ProbOver35K, prop.ProbOver45K, prop.ProbOver55K, prop.ProbOver65K,
		prop.ProbOver75K, prop.ProbOver85K, prop.ProbOver95K, prop.ProbOver105K, prop.ProbOver115K, prop.ProbOver125K,
		prop.Over25KLower95, prop.Over25KUpper95, prop.Over35KLower95, prop.Over35KUpper95,
		prop.Over45KLower95, prop.Over45KUpper95, prop.Over55KLower95, prop.Over55KUpper95,
		prop.Over65KLower95, prop.Over65KUpper95, prop.Over75KLower95, prop.Over75KUpper95,
		prop.Over85KLower95, prop.Over85KUpper95, prop.Over95KLower95, prop.Over95KUpper95,
		prop.Over105KLower95, prop.Over105KUpper95, prop.Over115KLower95, prop.Over115KUpper95,
		prop.Over125KLower95, prop.Over125KUpper95,
		prop.AvgStrikeouts, prop.IqrStrikeouts, prop.Q80Strikeouts, prop.Team, prop.StrikeoutPct, prop.WalkPct,
		prop.SwingingStrPct, prop.IP, prop.TotalStrikeouts, prop.TotalWalks, prop.TotalPlateAppearances,
		prop.TotalSwingingStrikes, prop.TotalPitches,
	}

	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	updates := []string{}
	for _, col := range columns[1:] {
		updates = append(updates, fmt.Sprintf("%s = EXCLUDED.%s", col, col))
	}

	sqlStr := fmt.Sprintf(
		`INSERT INTO pitcher_props (%s) VALUES (%s)
		ON CONFLICT (pitcherid, gamepk) DO UPDATE SET %s`,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(updates, ", "),
	)
	_, err := db.Exec(sqlStr, values...)
	return err
}
