package sim

import (
 //"github.com/logananthony/go-baseball/pkg/sim"
 "github.com/logananthony/go-baseball/pkg/models"
 "github.com/logananthony/go-baseball/pkg/fetcher"
 "github.com/logananthony/go-baseball/pkg/config"
 "github.com/logananthony/go-baseball/pkg/utils"

)

func SimulateAtBat(in []models.AtBatData ) []models.AtBatResult {

    db := config.ConnectDB()
    defer db.Close()

    balls := 0
    strikes := 0
    pitch_count := 0


    batterStands := fetcher.FetchBatterInfo(db, in[0].BatterId, in[0].GameYear)
    pitcherThrows := fetcher.FetchPitcherInfo(db, in[0].PitcherId, in[0].GameYear)

    if batterStands == "B" && pitcherThrows == "R" {
      batterStands = "L"
      } else if batterStands == "B" && pitcherThrows == "L" {
      batterStands = "R"
    }

    pitch_type_sequence := []string {}
    plate_x_sequence := []float64 {}
    plate_z_sequence := []float64 {}
    velocity_sequence := []float64 {}
    strike_sequence := []int {}
    ball_sequence := []int {}
    is_strike_sequence := []bool {}


    for {
    
      pitch_count += 1
      
      pitcher_freqs := fetcher.FetchPitcherFrequencies(db, in[0].PitcherId, batterStands)
      pitch_type_result := SimulatePitchType(pitcher_freqs, balls, strikes)
      pitch_covariance := fetcher.FetchPitcherCovarianceMean(db, int64(in[0].PitcherId), int64(in[0].GameYear))
      location_velo_result := SimulatePitchLocationVelo(pitch_covariance, pitch_type_result, batterStands, balls, strikes)
      is_strike_result := utils.IsPitchStrike(location_velo_result[0], location_velo_result[1])


      pitch_type_sequence = append(pitch_type_sequence, pitch_type_result)
      plate_x_sequence = append(plate_x_sequence, location_velo_result[0])
      plate_z_sequence = append(plate_z_sequence, location_velo_result[1])
      velocity_sequence = append(velocity_sequence, location_velo_result[2])
      strike_sequence = append(strike_sequence, strikes)
      ball_sequence = append(ball_sequence, balls)
      is_strike_sequence = append(is_strike_sequence, is_strike_result)

      if is_strike_result {
         strikes += 1
      } else {
         balls += 1
      }

      if strikes == 3 || balls == 4 {

        return []models.AtBatResult{{
          GameYear: in[0].GameYear,
          PitcherId: in[0].PitcherId,
          BatterId: in[0].BatterId,
          Strikes: strike_sequence, 
          Balls: ball_sequence, 
          PitchType: pitch_type_sequence,
          PlateX: plate_x_sequence,
          PlateZ: plate_z_sequence, 
          Velocity: velocity_sequence,
          IsStrike: is_strike_sequence,

        }}

      }

    }

}
