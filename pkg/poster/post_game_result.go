package poster

import (
	"database/sql"
	"fmt"

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

func minInt(vals ...int) int {
	if len(vals) == 0 {
		return 0
	}
	min := vals[0]
	for _, v := range vals[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

func InsertGameResult(db *sql.DB, gameId string, gameYear int, result models.GameResult) error {
	for _, pa := range result.PitchData {
		numPitches := minInt(
			len(pa.PitcherId),
			len(pa.BatterId),
			len(pa.Strikes),
			len(pa.Balls),
			len(pa.PitchCount),
			len(pa.PitchType),
			len(pa.PlateX),
			len(pa.PlateZ),
			len(pa.Zone),
			len(pa.Velocity),
			len(pa.IsStrike),
			len(pa.IsSwing),
		)

		for i := 0; i < numPitches; i++ {
			var isContact interface{}
			if pa.IsSwing[i] && i < len(pa.IsContact) {
				isContact = pa.IsContact[i]
			} else {
				isContact = nil
			}

			var eventType interface{}
			if i == numPitches-1 && len(pa.EventType) > 0 {
				eventType = pa.EventType[len(pa.EventType)-1]
			} else {
				eventType = nil
			}

      var exitVelocity any
      if i == numPitches-1 && len(pa.ExitVelocity) > 0 && i < len(pa.IsContact) && pa.IsContact[i] == "ball_in_play" {
        exitVelocity = pa.ExitVelocity[0]
      } else {
        exitVelocity = nil
      }

      var launchAngle any
      if i == numPitches-1 && len(pa.LaunchAngle) > 0 && i < len(pa.IsContact) && pa.IsContact[i] == "ball_in_play" {
        launchAngle = pa.LaunchAngle[0]
      } else {
        launchAngle = nil
      }

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
                event_type, exit_velocity, launch_angle
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
                $31, $32, $33
              )
          `,
				gameId, gameYear,
				result.AtBatNumber, result.Inning, result.InningTopBot,
				result.Outs, result.On1b, result.On2b, result.On3b,
				result.AwayScore, result.HomeScore,
				pa.PitcherId[i], toNullableString(pa.PitcherFullName, i), toNullableInt(pa.PitcherGameYear, i),
				pa.BatterId[i], toNullableString(pa.BatterFullName, i), toNullableInt(pa.BatterGameYear, i),
				toNullableString(pa.BatterStands, i), toNullableString(pa.PitcherThrows, i),
				pa.Strikes[i], pa.Balls[i], pa.PitchCount[i],
				pa.PitchType[i], pa.PlateX[i], pa.PlateZ[i], pa.Zone[i],
				pa.Velocity[i], pa.IsStrike[i], pa.IsSwing[i], isContact,
				eventType, exitVelocity, launchAngle,
			)

			if err != nil {
				fmt.Println("🚨 INSERT FAILED for AtBatNumber", result.AtBatNumber, ":", err)
				return err
			}

			fmt.Println("✅ Inserted pitch", i+1, "for AtBatNumber:", result.AtBatNumber)
		}
	}
	return nil
}

