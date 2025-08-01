package poster

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/logananthony/go-baseball/pkg/models"
)

// -- NULL HELPERS -- //
func toNullString(slice []string, i int) sql.NullString {
	if slice == nil || len(slice) <= i || slice[i] == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: slice[i], Valid: true}
}
func toNullInt(slice []int, i int) sql.NullInt64 {
	if slice == nil || len(slice) <= i {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(slice[i]), Valid: true}
}
func toNullFloat(slice []float64, i int) sql.NullFloat64 {
	if slice == nil || len(slice) <= i {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: slice[i], Valid: true}
}

// -- BULK INSERT FUNCTION -- //
func InsertGameResult(db *sql.DB, gameId string, gamepk int, jobId string, gameYear int, result models.GameResult) error {
	pa := result.PAResult
	numPitches := len(pa.PitcherId)
	batchSize := 100

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	columns := []string{
		"game_id", "game_year",
		"at_bat_number", "inning", "inning_topbot",
		"outs", "on1B", "on2B", "on3B",
		"runner_on1b_id", "runner_on2b_id", "runner_on3b_id",
		"rbi_batter_id", "rbi_count",
		"away_score", "home_score",
		"pitcherid", "pitcher_fullname", "pitcher_game_year",
		"batterid", "batter_fullname", "batter_game_year",
		"batter_stands", "pitcher_throws",
		"strikes", "balls", "pitch_count",
		"pitch_type", "plate_x", "plate_z", "zone",
		"velocity", "is_strike", "is_swing", "is_contact",
		"event_type", "exit_velocity", "launch_angle", "spray_angle",
		"created_at", "jobId", "gamepk",
	}
	placeholdersPerRow := len(columns)
	totalBatches := (numPitches + batchSize - 1) / batchSize

	for b := 0; b < totalBatches; b++ {
		start := b * batchSize
		end := start + batchSize
		if end > numPitches {
			end = numPitches
		}

		var valueStrings []string
		var args []interface{}
		argIdx := 1

		for i := start; i < end; i++ {
			valueStrings = append(valueStrings, "("+makePlaceholders(placeholdersPerRow, argIdx)+")")
			argIdx += placeholdersPerRow

			args = append(args,
				gameId, gameYear,
				pa.AtBatNumber[i], pa.Inning[i], pa.InningTopBot[i],
				pa.Outs[i], pa.On1b[i], pa.On2b[i], pa.On3b[i],
				toNullInt(pa.RunnerOn1bID, i), toNullInt(pa.RunnerOn2bID, i), toNullInt(pa.RunnerOn3bID, i),
				toNullInt(pa.RbiBatterID, i), toNullInt(pa.RbiCount, i),
				pa.AwayScore[i], pa.HomeScore[i],
				pa.PitcherId[i], toNullString(pa.PitcherFullName, i), toNullInt(pa.PitcherGameYear, i),
				pa.BatterId[i], toNullString(pa.BatterFullName, i), toNullInt(pa.BatterGameYear, i),
				toNullString(pa.BatterStands, i), toNullString(pa.PitcherThrows, i),
				pa.Strikes[i], pa.Balls[i], pa.PitchCount[i],
				toNullString(pa.PitchType, i), pa.PlateX[i], pa.PlateZ[i], pa.Zone[i],
				pa.Velocity[i], pa.IsStrike[i], pa.IsSwing[i], toNullString(pa.IsContact, i),
				toNullString(pa.EventType, i),
				toNullFloat(pa.ExitVelocity, i), toNullFloat(pa.LaunchAngle, i), toNullFloat(pa.SprayAngle, i),
				time.Now(), jobId, gamepk,
			)
		}

		sqlStr := fmt.Sprintf(
			"INSERT INTO game_result (%s) VALUES %s",
			strings.Join(columns, ", "),
			strings.Join(valueStrings, ","),
		)

		_, err := tx.Exec(sqlStr, args...)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("batch insert failed (rows %d-%d): %w", start, end-1, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	fmt.Printf("✅ Inserted %d pitches for gamePk %d (batched)\n", numPitches, gamepk)
	return nil
}

// -- Placeholder generator -- //
func makePlaceholders(n, offset int) string {
	vals := make([]string, n)
	for i := 0; i < n; i++ {
		vals[i] = fmt.Sprintf("$%d", offset+i)
	}
	return strings.Join(vals, ",")
}
