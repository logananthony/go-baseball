package poster

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/logananthony/go-baseball/pkg/models"
)

// InsertGameResultAggCore inserts or upserts a single aggregated row
func InsertGameResultAggCore(db *sql.DB, agg models.GameResultAggCore) error {
	columns := []string{
		"gamepk", "total_over_35", "total_over_45", "total_over_65", "total_over_75", "total_over_85", "total_over_95", "total_over_105",
		"spread_minus_25", "spread_minus_15", "spread_plus_15", "spread_plus_25",
		"moneyline_home_win", "moneyline_away_win", "gamedate", "hometeamabbr", "awayteamabbr", "homepitchername", "awaypitchername",
		"total_over_15", "total_over_25", "total_over_55", "total_over_115", "total_over_125",
		"spread_minus_55", "spread_minus_45", "spread_minus_35", "spread_plus_35", "spread_plus_45", "spread_plus_55",
		"home_total_over_05", "home_total_over_15", "home_total_over_25", "home_total_over_35", "home_total_over_45", "home_total_over_55", "home_total_over_65",
		"away_total_over_05", "away_total_over_15", "away_total_over_25", "away_total_over_35", "away_total_over_45", "away_total_over_55", "away_total_over_65",
		"std_total_runs", "std_home_score", "std_away_score", "ml_var", "std_spread",
	}
	values := []interface{}{
		agg.GamePk, agg.TotalOver35, agg.TotalOver45, agg.TotalOver65, agg.TotalOver75, agg.TotalOver85, agg.TotalOver95, agg.TotalOver105,
		agg.SpreadMinus25, agg.SpreadMinus15, agg.SpreadPlus15, agg.SpreadPlus25,
		agg.MoneylineHomeWin, agg.MoneylineAwayWin, agg.GameDate, agg.HomeTeamAbbr, agg.AwayTeamAbbr, agg.HomePitcherName, agg.AwayPitcherName,
		agg.TotalOver15, agg.TotalOver25, agg.TotalOver55, agg.TotalOver115, agg.TotalOver125,
		agg.SpreadMinus55, agg.SpreadMinus45, agg.SpreadMinus35, agg.SpreadPlus35, agg.SpreadPlus45, agg.SpreadPlus55,
		agg.HomeTotalOver05, agg.HomeTotalOver15, agg.HomeTotalOver25, agg.HomeTotalOver35, agg.HomeTotalOver45, agg.HomeTotalOver55, agg.HomeTotalOver65,
		agg.AwayTotalOver05, agg.AwayTotalOver15, agg.AwayTotalOver25, agg.AwayTotalOver35, agg.AwayTotalOver45, agg.AwayTotalOver55, agg.AwayTotalOver65,
		agg.StdTotalRuns, agg.StdHomeScore, agg.StdAwayScore, agg.MlVar, agg.StdSpread,
	}
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	updates := []string{}
	for _, col := range columns[1:] { // skip "gamepk" as primary key
		updates = append(updates, fmt.Sprintf("%s = EXCLUDED.%s", col, col))
	}

	sqlStr := fmt.Sprintf(
		`INSERT INTO game_result_agg_core (%s) VALUES (%s)
		 ON CONFLICT (gamepk) DO UPDATE SET %s`,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(updates, ", "),
	)

	_, err := db.Exec(sqlStr, values...)
	return err
}
