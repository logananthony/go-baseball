package poster

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/logananthony/go-baseball/pkg/models"
)

func toNullableString(slice []string, i int) string {
	if slice == nil || len(slice) <= i {
		return ""
	}
	return slice[i]
}

func toNullableInt(slice []int, i int) int {
	if slice == nil || len(slice) <= i {
		return 0
	}
	return slice[i]
}

func InsertGameResult(db *sql.DB, gameId string, gameYear int, result models.GameResult) error {
	pa := result.PAResult

	// Ensure all slices in PlateAppearanceResult are of the same length
	numPitches := len(pa.PitcherId) // Use one of the slices to determine the number of pitches
	for i := 0; i < numPitches; i++ {
		_, err := db.Exec(`
            INSERT INTO game_result (
                game_id, game_year,
                at_bat_number, inning, inning_topbot,
                outs, on1B, on2B, on3B,
                away_score, home_score,
                pitcherid, pitcher_fullname, pitcher_game_year,
                batterid, batter_fullname, batter_game_year,
                batter_stands, pitcher_throws,
                strikes, balls, pitch_count,
                pitch_type, plate_x, plate_z, zone,
                velocity, is_strike, is_swing, is_contact,
                event_type, exit_velocity, launch_angle, spray_angle,
                created_at
            ) VALUES (
                $1, $2, $3, $4, $5,
                $6, $7, $8, $9,
                $10, $11,
                $12, $13, $14,
                $15, $16, $17,
                $18, $19,
                $20, $21, $22,
                $23, $24, $25, $26,
                $27, $28, $29, $30,
                $31, $32, $33, $34,
                $35
            )
        `,
			gameId, gameYear,
			pa.AtBatNumber[i], pa.Inning[i], pa.InningTopBot[i],
			pa.Outs[i], pa.On1b[i], pa.On2b[i], pa.On3b[i],
			pa.AwayScore[i], pa.HomeScore[i],
			pa.PitcherId[i], toNullableString(pa.PitcherFullName, i), toNullableInt(pa.PitcherGameYear, i),
			pa.BatterId[i], toNullableString(pa.BatterFullName, i), toNullableInt(pa.BatterGameYear, i),
			toNullableString(pa.BatterStands, i), toNullableString(pa.PitcherThrows, i),
			pa.Strikes[i], pa.Balls[i], pa.PitchCount[i],
			toNullableString(pa.PitchType, i), pa.PlateX[i], pa.PlateZ[i], pa.Zone[i],
			pa.Velocity[i], pa.IsStrike[i], pa.IsSwing[i], toNullableString(pa.IsContact, i),
			toNullableString(pa.EventType, i), pa.ExitVelocity[i], pa.LaunchAngle[i], pa.SprayAngle[i], time.Now(),
		)

		if err != nil {
			fmt.Printf("🚨 INSERT FAILED for AtBatNumber %d, Pitch %d: %v\n", pa.AtBatNumber[i], i+1, err)
			return err
		}

		fmt.Printf("✅ Inserted pitch %d for AtBatNumber: %d\n", i+1, pa.AtBatNumber[i])
	}
	return nil
}
