package poster

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/logananthony/go-baseball/pkg/models"
)

func InsertBatterProps(db *sql.DB, prop models.BatterProps) error {
	columns := []string{
		"batterid", "gamepk", "battername", "battingorder", "team", "gamedate", "num_simulations",
		"prob_over_0_5_hits", "prob_over_1_5_hits", "over_0_5_hits_lower95", "over_0_5_hits_upper95",
		"over_1_5_hits_lower95", "over_1_5_hits_upper95", "avg_hits", "iqr_hits", "q80_hits",
		"prob_over_0_5_singles", "prob_over_1_5_singles", "over_0_5_singles_lower95", "over_0_5_singles_upper95",
		"over_1_5_singles_lower95", "over_1_5_singles_upper95", "avg_singles", "iqr_singles", "q80_singles",
		"prob_over_0_5_doubles", "prob_over_1_5_doubles", "over_0_5_doubles_lower95", "over_0_5_doubles_upper95",
		"over_1_5_doubles_lower95", "over_1_5_doubles_upper95", "avg_doubles", "iqr_doubles", "q80_doubles",
		"prob_over_0_5_triples", "prob_over_1_5_triples", "over_0_5_triples_lower95", "over_0_5_triples_upper95",
		"over_1_5_triples_lower95", "over_1_5_triples_upper95", "avg_triples", "iqr_triples", "q80_triples",
		"prob_over_0_5_homeruns", "prob_over_1_5_homeruns", "over_0_5_homeruns_lower95", "over_0_5_homeruns_upper95",
		"over_1_5_homeruns_lower95", "over_1_5_homeruns_upper95", "avg_homeruns", "iqr_homeruns", "q80_homeruns",
		"batting_avg", "slugging_pct", "on_base_pct", "strikeout_pct", "walk_pct",
	}

	values := []interface{}{
		prop.BatterID, prop.GamePk, prop.BatterName, prop.BattingOrder, prop.Team, prop.GameDate, prop.NumSimulations,
		prop.ProbOver05Hits, prop.ProbOver15Hits, prop.Over05HitsLower95, prop.Over05HitsUpper95,
		prop.Over15HitsLower95, prop.Over15HitsUpper95, prop.AvgHits, prop.IqrHits, prop.Q80Hits,
		prop.ProbOver05Singles, prop.ProbOver15Singles, prop.Over05SinglesLower95, prop.Over05SinglesUpper95,
		prop.Over15SinglesLower95, prop.Over15SinglesUpper95, prop.AvgSingles, prop.IqrSingles, prop.Q80Singles,
		prop.ProbOver05Doubles, prop.ProbOver15Doubles, prop.Over05DoublesLower95, prop.Over05DoublesUpper95,
		prop.Over15DoublesLower95, prop.Over15DoublesUpper95, prop.AvgDoubles, prop.IqrDoubles, prop.Q80Doubles,
		prop.ProbOver05Triples, prop.ProbOver15Triples, prop.Over05TriplesLower95, prop.Over05TriplesUpper95,
		prop.Over15TriplesLower95, prop.Over15TriplesUpper95, prop.AvgTriples, prop.IqrTriples, prop.Q80Triples,
		prop.ProbOver05Homeruns, prop.ProbOver15Homeruns, prop.Over05HomerunsLower95, prop.Over05HomerunsUpper95,
		prop.Over15HomerunsLower95, prop.Over15HomerunsUpper95, prop.AvgHomeruns, prop.IqrHomeruns, prop.Q80Homeruns,
		prop.BattingAvg, prop.SluggingPct, prop.OnBasePct, prop.StrikeoutPct, prop.WalkPct,
	}

	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	updates := []string{}
	for _, col := range columns[1:] { // skip PK
		updates = append(updates, fmt.Sprintf("%s = EXCLUDED.%s", col, col))
	}

	sqlStr := fmt.Sprintf(
		`INSERT INTO batter_props (%s) VALUES (%s)
		ON CONFLICT (batterid, gamepk) DO UPDATE SET %s`,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(updates, ", "),
	)
	_, err := db.Exec(sqlStr, values...)
	return err
}
