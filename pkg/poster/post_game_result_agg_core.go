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
		"gamepk",
		// Game totals (point + CI)
		"total_over_15", "total_over_15_lower95", "total_over_15_upper95",
		"total_over_25", "total_over_25_lower95", "total_over_25_upper95",
		"total_over_35", "total_over_35_lower95", "total_over_35_upper95",
		"total_over_45", "total_over_45_lower95", "total_over_45_upper95",
		"total_over_55", "total_over_55_lower95", "total_over_55_upper95",
		"total_over_65", "total_over_65_lower95", "total_over_65_upper95",
		"total_over_75", "total_over_75_lower95", "total_over_75_upper95",
		"total_over_85", "total_over_85_lower95", "total_over_85_upper95",
		"total_over_95", "total_over_95_lower95", "total_over_95_upper95",
		"total_over_105", "total_over_105_lower95", "total_over_105_upper95",
		"total_over_115", "total_over_115_lower95", "total_over_115_upper95",
		"total_over_125", "total_over_125_lower95", "total_over_125_upper95",

		// Home team totals (point + CI)
		"home_total_over_05", "home_total_over_05_lower95", "home_total_over_05_upper95",
		"home_total_over_15", "home_total_over_15_lower95", "home_total_over_15_upper95",
		"home_total_over_25", "home_total_over_25_lower95", "home_total_over_25_upper95",
		"home_total_over_35", "home_total_over_35_lower95", "home_total_over_35_upper95",
		"home_total_over_45", "home_total_over_45_lower95", "home_total_over_45_upper95",
		"home_total_over_55", "home_total_over_55_lower95", "home_total_over_55_upper95",
		"home_total_over_65", "home_total_over_65_lower95", "home_total_over_65_upper95",

		// Away team totals (point + CI)
		"away_total_over_05", "away_total_over_05_lower95", "away_total_over_05_upper95",
		"away_total_over_15", "away_total_over_15_lower95", "away_total_over_15_upper95",
		"away_total_over_25", "away_total_over_25_lower95", "away_total_over_25_upper95",
		"away_total_over_35", "away_total_over_35_lower95", "away_total_over_35_upper95",
		"away_total_over_45", "away_total_over_45_lower95", "away_total_over_45_upper95",
		"away_total_over_55", "away_total_over_55_lower95", "away_total_over_55_upper95",
		"away_total_over_65", "away_total_over_65_lower95", "away_total_over_65_upper95",

		// Spreads (point + CI)
		"spread_minus_55", "spread_minus_55_lower95", "spread_minus_55_upper95",
		"spread_minus_45", "spread_minus_45_lower95", "spread_minus_45_upper95",
		"spread_minus_35", "spread_minus_35_lower95", "spread_minus_35_upper95",
		"spread_minus_25", "spread_minus_25_lower95", "spread_minus_25_upper95",
		"spread_minus_15", "spread_minus_15_lower95", "spread_minus_15_upper95",
		"spread_plus_15", "spread_plus_15_lower95", "spread_plus_15_upper95",
		"spread_plus_25", "spread_plus_25_lower95", "spread_plus_25_upper95",
		"spread_plus_35", "spread_plus_35_lower95", "spread_plus_35_upper95",
		"spread_plus_45", "spread_plus_45_lower95", "spread_plus_45_upper95",
		"spread_plus_55", "spread_plus_55_lower95", "spread_plus_55_upper95",

		"moneyline_home_win", "ml_home_win_lower95", "ml_home_win_upper95",
		"moneyline_away_win",
		"gamedate", "hometeamabbr", "awayteamabbr", "homepitchername", "awaypitchername",

		"std_total_runs", "iqr_total_runs", "q80_total_runs",
		"std_home_score", "iqr_home_score", "q80_home_score", "home_score_lower95", "home_score_upper95",
		"std_away_score", "iqr_away_score", "q80_away_score", "away_score_lower95", "away_score_upper95",
		"ml_var",
		"std_spread", "iqr_spread", "q80_spread", "spread_lower95", "spread_upper95",
	}

	values := []interface{}{
		agg.GamePk,

		// Game totals (point + CI)
		agg.TotalOver15, agg.TotalOver15Lower95, agg.TotalOver15Upper95,
		agg.TotalOver25, agg.TotalOver25Lower95, agg.TotalOver25Upper95,
		agg.TotalOver35, agg.TotalOver35Lower95, agg.TotalOver35Upper95,
		agg.TotalOver45, agg.TotalOver45Lower95, agg.TotalOver45Upper95,
		agg.TotalOver55, agg.TotalOver55Lower95, agg.TotalOver55Upper95,
		agg.TotalOver65, agg.TotalOver65Lower95, agg.TotalOver65Upper95,
		agg.TotalOver75, agg.TotalOver75Lower95, agg.TotalOver75Upper95,
		agg.TotalOver85, agg.TotalOver85Lower95, agg.TotalOver85Upper95,
		agg.TotalOver95, agg.TotalOver95Lower95, agg.TotalOver95Upper95,
		agg.TotalOver105, agg.TotalOver105Lower95, agg.TotalOver105Upper95,
		agg.TotalOver115, agg.TotalOver115Lower95, agg.TotalOver115Upper95,
		agg.TotalOver125, agg.TotalOver125Lower95, agg.TotalOver125Upper95,

		// Home team totals (point + CI)
		agg.HomeTotalOver05, agg.HomeTotalOver05Lower95, agg.HomeTotalOver05Upper95,
		agg.HomeTotalOver15, agg.HomeTotalOver15Lower95, agg.HomeTotalOver15Upper95,
		agg.HomeTotalOver25, agg.HomeTotalOver25Lower95, agg.HomeTotalOver25Upper95,
		agg.HomeTotalOver35, agg.HomeTotalOver35Lower95, agg.HomeTotalOver35Upper95,
		agg.HomeTotalOver45, agg.HomeTotalOver45Lower95, agg.HomeTotalOver45Upper95,
		agg.HomeTotalOver55, agg.HomeTotalOver55Lower95, agg.HomeTotalOver55Upper95,
		agg.HomeTotalOver65, agg.HomeTotalOver65Lower95, agg.HomeTotalOver65Upper95,

		// Away team totals (point + CI)
		agg.AwayTotalOver05, agg.AwayTotalOver05Lower95, agg.AwayTotalOver05Upper95,
		agg.AwayTotalOver15, agg.AwayTotalOver15Lower95, agg.AwayTotalOver15Upper95,
		agg.AwayTotalOver25, agg.AwayTotalOver25Lower95, agg.AwayTotalOver25Upper95,
		agg.AwayTotalOver35, agg.AwayTotalOver35Lower95, agg.AwayTotalOver35Upper95,
		agg.AwayTotalOver45, agg.AwayTotalOver45Lower95, agg.AwayTotalOver45Upper95,
		agg.AwayTotalOver55, agg.AwayTotalOver55Lower95, agg.AwayTotalOver55Upper95,
		agg.AwayTotalOver65, agg.AwayTotalOver65Lower95, agg.AwayTotalOver65Upper95,

		// Spreads (point + CI)
		agg.SpreadMinus55, agg.SpreadMinus55Lower95, agg.SpreadMinus55Upper95,
		agg.SpreadMinus45, agg.SpreadMinus45Lower95, agg.SpreadMinus45Upper95,
		agg.SpreadMinus35, agg.SpreadMinus35Lower95, agg.SpreadMinus35Upper95,
		agg.SpreadMinus25, agg.SpreadMinus25Lower95, agg.SpreadMinus25Upper95,
		agg.SpreadMinus15, agg.SpreadMinus15Lower95, agg.SpreadMinus15Upper95,
		agg.SpreadPlus15, agg.SpreadPlus15Lower95, agg.SpreadPlus15Upper95,
		agg.SpreadPlus25, agg.SpreadPlus25Lower95, agg.SpreadPlus25Upper95,
		agg.SpreadPlus35, agg.SpreadPlus35Lower95, agg.SpreadPlus35Upper95,
		agg.SpreadPlus45, agg.SpreadPlus45Lower95, agg.SpreadPlus45Upper95,
		agg.SpreadPlus55, agg.SpreadPlus55Lower95, agg.SpreadPlus55Upper95,

		agg.MoneylineHomeWin, agg.MlHomeWinLower95, agg.MlHomeWinUpper95,
		agg.MoneylineAwayWin,
		agg.GameDate, agg.HomeTeamAbbr, agg.AwayTeamAbbr, agg.HomePitcherName, agg.AwayPitcherName,

		agg.StdTotalRuns, agg.IqrTotalRuns, agg.Q80TotalRuns,
		agg.StdHomeScore, agg.IqrHomeScore, agg.Q80HomeScore, agg.HomeScoreLower95, agg.HomeScoreUpper95,
		agg.StdAwayScore, agg.IqrAwayScore, agg.Q80AwayScore, agg.AwayScoreLower95, agg.AwayScoreUpper95,
		agg.MlVar,
		agg.StdSpread, agg.IqrSpread, agg.Q80Spread, agg.SpreadLower95, agg.SpreadUpper95,
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
