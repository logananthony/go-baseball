package sim

import (
 //"github.com/logananthony/go-baseball/pkg/sim"
"github.com/logananthony/go-baseball/pkg/models"
 "github.com/logananthony/go-baseball/pkg/fetcher"
 "github.com/logananthony/go-baseball/pkg/config"
 "github.com/logananthony/go-baseball/pkg/utils"
 "fmt"

)

func SimulateAtBat(in []models.PlateAppearanceData) []models.PlateAppearanceResult {

    db := config.ConnectDB()
    defer db.Close()

    balls := 0
    strikes := 0
    pitch_count := 0


    batterStands := fetcher.FetchBatterInfo(db, in[0].BatterId, in[0].GameYear)
    pitcherThrows := fetcher.FetchPitcherInfo(db, in[0].PitcherId, in[0].GameYear)
    
    batterSwingProbs, _ := fetcher.FetchBatterSwingPercentage(db, in[0].BatterId, in[0].GameYear)
    batterContactProbs, _ := fetcher.FetchBatterContactPercentage(db, in[0].BatterId, in[0].GameYear)

    batterHitProbs, err := fetcher.FetchBatterHitType(db, 660688, 2024)
    if err != nil {
        fmt.Println("Fetcher error:", err)
    }


    if batterStands == "B" && pitcherThrows == "R" {
      batterStands = "L"
      } else if batterStands == "B" && pitcherThrows == "L" {
      batterStands = "R"
    }
   

    game_year_sequence := []int {}
    batterid_sequence := []int {}
    pitcherid_sequence := []int {}
  
    pitch_type_sequence := []string {}
    plate_x_sequence := []float64 {}
    plate_z_sequence := []float64 {}
    velocity_sequence := []float64 {}
    strike_sequence := []int {}
    ball_sequence := []int {}
    pitch_count_sequence := []int {}
    is_strike_sequence := []bool {}
    is_swing_sequence := []bool {}
    is_contact_sequence := []string {}
    event_type_sequence := []string {}


    for {
    
      pitch_count += 1
      
      
      pitcher_freqs := fetcher.FetchPitcherFrequencies(db, in[0].PitcherId, batterStands)
      pitch_type_result := SimulatePitchType(pitcher_freqs, balls, strikes)
      pitch_covariance := fetcher.FetchPitcherCovarianceMean(db, int64(in[0].PitcherId), int64(in[0].GameYear))
      location_velo_result := SimulatePitchLocationVelo(pitch_covariance, pitch_type_result, batterStands, balls, strikes)
      is_strike_result := utils.IsPitchStrike(location_velo_result[0], location_velo_result[1])
      is_swing_result := SimulateSwingDecision(batterSwingProbs, batterStands, pitcherThrows, pitch_type_result, location_velo_result[0], location_velo_result[1])
      is_contact_result := SimulateContactPercentage(batterContactProbs, batterStands, pitcherThrows, pitch_type_result, location_velo_result[0], location_velo_result[1])
      event_type_result := SimulateBatterHitType(batterHitProbs, batterStands, pitcherThrows, pitch_type_result, location_velo_result[0], location_velo_result[1], 
                                                                                                                                  location_velo_result[2])
      game_year_sequence = append(game_year_sequence, in[0].GameYear)
      batterid_sequence = append(batterid_sequence, in[0].BatterId)
      pitcherid_sequence = append(pitcherid_sequence, in[0].PitcherId)
      pitch_type_sequence = append(pitch_type_sequence, pitch_type_result)
      plate_x_sequence = append(plate_x_sequence, location_velo_result[0])
      plate_z_sequence = append(plate_z_sequence, location_velo_result[1])
      velocity_sequence = append(velocity_sequence, location_velo_result[2])
      is_strike_sequence = append(is_strike_sequence, is_strike_result)
      is_swing_sequence = append(is_swing_sequence, is_swing_result)
      strike_sequence = append(strike_sequence, strikes)
      ball_sequence = append(ball_sequence, balls)
      pitch_count_sequence = append(pitch_count_sequence, pitch_count)

      if is_swing_result {
          // Batter swung
          is_contact_sequence = append(is_contact_sequence, is_contact_result)
          
          switch is_contact_result {
          case "swinging_strike":
              strikes += 1
          case "foul":
              if strikes < 2 {
                  strikes += 1 // foul with less than 2 strikes
              }
          case "ball_in_play":
                event_type_sequence = append(event_type_sequence, event_type_result)
                return []models.PlateAppearanceResult{{
                  GameYear: game_year_sequence,
                  PitcherId: pitcherid_sequence,
                  BatterId: batterid_sequence,
                  Strikes: strike_sequence, 
                  Balls: ball_sequence, 
                  PitchCount: pitch_count_sequence,
                  PitchType: pitch_type_sequence,
                  PlateX: plate_x_sequence,
                  PlateZ: plate_z_sequence, 
                  Velocity: velocity_sequence,
                  IsStrike: is_strike_sequence,
                  IsSwing: is_swing_sequence,
                  IsContact: is_contact_sequence,
                  EventType: event_type_sequence,

                }}
                



              // Let it resolve in your end condition (put in play)
          }

      } else {
          // Batter didn't swing
          if is_strike_result {
              strikes += 1 // called strike
          } else {
              balls += 1 // ball taken
          }
      }

      if strikes == 3 || balls == 4 {

      switch {

      case strikes == 3:

            event_type_sequence = append(event_type_sequence, "strikeout")
            return []models.PlateAppearanceResult{{
              GameYear: game_year_sequence,
              PitcherId: pitcherid_sequence,
              BatterId: batterid_sequence,
              Strikes: strike_sequence, 
              Balls: ball_sequence, 
              PitchCount: pitch_count_sequence,
              PitchType: pitch_type_sequence,
              PlateX: plate_x_sequence,
              PlateZ: plate_z_sequence, 
              Velocity: velocity_sequence,
              IsStrike: is_strike_sequence,
              IsSwing: is_swing_sequence,
              IsContact: is_contact_sequence,
              EventType: event_type_sequence,

            }}
      case balls == 4:
          event_type_sequence = append(event_type_sequence, "walk")
            return []models.PlateAppearanceResult{{
              GameYear: game_year_sequence,
              PitcherId: pitcherid_sequence,
              BatterId: batterid_sequence,
              Strikes: strike_sequence, 
              Balls: ball_sequence,
              PitchCount: pitch_count_sequence,
              PitchType: pitch_type_sequence,
              PlateX: plate_x_sequence,
              PlateZ: plate_z_sequence, 
              Velocity: velocity_sequence,
              IsStrike: is_strike_sequence,
              IsSwing: is_swing_sequence,
              IsContact: is_contact_sequence,
              EventType: event_type_sequence,

          }}

      }

    }

  }  

}
