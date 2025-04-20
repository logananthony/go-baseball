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

    //fmt.Printf("Fetching ID: %d | Season: %d\n", in[0].BatterId, in[0].GameYear)

    batterInfo, err := fetcher.FetchPlayerInfo(db, &in[0].BatterId, &in[0].BatterGameYear)
    if err != nil {
        // handle error, log.Fatal, return, etc.
    }
    pitcherInfo, err := fetcher.FetchPlayerInfo(db, &in[0].PitcherId, &in[0].PitcherGameYear)
    if err != nil {
        // handle error
    }
  //fmt.Println(in[0].BatterId)


    if len(batterInfo) == 0 {
        panic(fmt.Sprintf("No batter info found for ID %d", in[0].BatterId))
    }

    if len(pitcherInfo) == 0 {
        panic(fmt.Sprintf("No pitcher info found for ID %d", in[0].PitcherId))
    }
  
   batterStands := batterInfo[0]
   pitcherThrows := pitcherInfo[0]

  if batterStands.BatSide != nil {
      switch *batterStands.BatSide {
      case "Left":
          left := "L"
          batterStands.BatSide = &left
      case "Right":
          right := "R"
          batterStands.BatSide = &right
      case "Switch":
          if pitcherThrows.PitchHand != nil {
              if *pitcherThrows.PitchHand == "Right" {
                  left := "L"
                  batterStands.BatSide = &left
              } else if *pitcherThrows.PitchHand == "Left" {
                  right := "R"
                  batterStands.BatSide = &right
              }
          }
      }
  }

  if pitcherThrows.PitchHand != nil {
      switch *pitcherThrows.PitchHand {
      case "Left":
          left := "L"
          pitcherThrows.PitchHand = &left
      case "Right":
          right := "R"
          pitcherThrows.PitchHand = &right
      }
  }

//fmt.Printf("batterStands: %+v\n", *batterStands.BatSide)

    batterSwingProbs, _ := fetcher.FetchBatterSwingPercentage(db, in[0].BatterId, in[0].BatterGameYear)
    batterContactProbs, _ := fetcher.FetchBatterContactPercentage(db, in[0].BatterId, in[0].BatterGameYear)

    batterHitProbs, err := fetcher.FetchBatterHitType(db, in[0].BatterId, in[0].BatterGameYear)
    if err != nil {
        fmt.Println("Fetcher error:", err)
    }

    pitcher_game_year_sequence := []int {}
    batter_game_year_sequence := []int {}
    pitcher_full_name_sequence := []string {}
    batter_full_name_sequence := []string {}
    batterid_sequence := []int {}
    pitcherid_sequence := []int {}
    stand_sequence := []string {}
    throws_sequence := []string {}
  
    pitch_type_sequence := []string {}
    plate_x_sequence := []float64 {}
    plate_z_sequence := []float64 {}
    zone_sequence := []int {}
    velocity_sequence := []float64 {}
    strike_sequence := []int {}
    ball_sequence := []int {}
    pitch_count_sequence := []int {}
    is_strike_sequence := []bool {}
    is_swing_sequence := []bool {}
    is_contact_sequence := []string {}
    exit_velocity_sequence := []float64 {}
    event_type_sequence := []string {}


    for {
    
      pitch_count += 1
      
      
      pitcher_freqs := fetcher.FetchPitcherFrequencies(db, in[0].PitcherId, *batterStands.BatSide)
      pitch_type_result := SimulatePitchType(pitcher_freqs, balls, strikes)
      pitch_covariance := fetcher.FetchPitcherCovarianceMean(db, int64(in[0].PitcherId), int64(in[0].PitcherGameYear))
      location_velo_result := SimulatePitchLocationVelo(pitch_covariance, pitch_type_result, *batterStands.BatSide, balls, strikes)
      zone_result := utils.GetPitchZone(location_velo_result[0], location_velo_result[1])
      is_strike_result := utils.IsPitchStrike(location_velo_result[0], location_velo_result[1])
      is_swing_result := SimulateSwingDecision(batterSwingProbs, *batterStands.BatSide, *pitcherThrows.PitchHand, pitch_type_result, location_velo_result[0], location_velo_result[1])
      is_contact_result := SimulateContactPercentage(batterContactProbs, *batterStands.BatSide, *pitcherThrows.PitchHand, pitch_type_result, location_velo_result[0], location_velo_result[1])
      event_type_result := SimulateBatterHitType(batterHitProbs, *batterStands.BatSide, *pitcherThrows.PitchHand, pitch_type_result, location_velo_result[0], location_velo_result[1], location_velo_result[2])

evRows := fetcher.FetchEVDistributions(
	db,
	in[0].BatterGameYear,
	in[0].BatterId,
	*batterStands.BatSide,
	*pitcherThrows.PitchHand,
	utils.StrToNull(event_type_result),
	utils.StrToNull(pitch_type_result),
	utils.IntToNull(zone_result), // If zone_result is -1, it'll convert to NULL
	utils.StrToNull(utils.GetVelocityBucket(location_velo_result[2])),
)


      agg := AggregateEVDistributions(evRows)
      ev_result := SampleFromAggregatedDistribution(agg)
    
      pitcher_game_year_sequence = append(pitcher_game_year_sequence, in[0].PitcherGameYear)
      batter_game_year_sequence = append(batter_game_year_sequence, in[0].BatterGameYear)
      pitcher_full_name_sequence = append(pitcher_full_name_sequence, *pitcherThrows.FullName)
      batter_full_name_sequence = append(batter_full_name_sequence, *batterStands.FullName)
      batterid_sequence = append(batterid_sequence, in[0].BatterId)
      pitcherid_sequence = append(pitcherid_sequence, in[0].PitcherId)
      stand_sequence = append(stand_sequence, *batterStands.BatSide)
      throws_sequence = append(throws_sequence, *pitcherThrows.PitchHand)
      pitch_type_sequence = append(pitch_type_sequence, pitch_type_result)
      plate_x_sequence = append(plate_x_sequence, location_velo_result[0])
      plate_z_sequence = append(plate_z_sequence, location_velo_result[1])
      zone_sequence = append(zone_sequence, zone_result)
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
                exit_velocity_sequence = append(exit_velocity_sequence, ev_result)
                return []models.PlateAppearanceResult{{
                  PitcherGameYear: pitcher_game_year_sequence,
                  PitcherFullName: pitcher_full_name_sequence,
                  PitcherId: pitcherid_sequence,
                  BatterGameYear: batter_game_year_sequence,
                  BatterFullName: batter_full_name_sequence,
                  BatterId: batterid_sequence,
                  BatterStands: stand_sequence,
                  PitcherThrows: throws_sequence,
                  Strikes: strike_sequence, 
                  Balls: ball_sequence, 
                  PitchCount: pitch_count_sequence,
                  PitchType: pitch_type_sequence,
                  PlateX: plate_x_sequence,
                  PlateZ: plate_z_sequence, 
                  Zone: zone_sequence,
                  Velocity: velocity_sequence,
                  IsStrike: is_strike_sequence,
                  IsSwing: is_swing_sequence,
                  IsContact: is_contact_sequence,
                  ExitVelocity: exit_velocity_sequence,
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
            exit_velocity_sequence = append(exit_velocity_sequence, 0)
            return []models.PlateAppearanceResult{{
                PitcherGameYear: pitcher_game_year_sequence,
                PitcherFullName: pitcher_full_name_sequence,
                PitcherId: pitcherid_sequence,
                BatterGameYear: batter_game_year_sequence,
                BatterFullName: batter_full_name_sequence,
                BatterId: batterid_sequence,
                BatterStands: stand_sequence,
                PitcherThrows: throws_sequence,
                Strikes: strike_sequence, 
                Balls: ball_sequence, 
                PitchCount: pitch_count_sequence,
                PitchType: pitch_type_sequence,
                PlateX: plate_x_sequence,
                PlateZ: plate_z_sequence, 
                Zone: zone_sequence,
                Velocity: velocity_sequence,
                IsStrike: is_strike_sequence,
                IsSwing: is_swing_sequence,
                IsContact: is_contact_sequence,
                ExitVelocity: exit_velocity_sequence,
                EventType: event_type_sequence,
            }}
      case balls == 4:
          event_type_sequence = append(event_type_sequence, "walk")
          exit_velocity_sequence = append(exit_velocity_sequence, 0)
            return []models.PlateAppearanceResult{{
                PitcherGameYear: pitcher_game_year_sequence,
                PitcherFullName: pitcher_full_name_sequence,
                PitcherId: pitcherid_sequence,
                BatterGameYear: batter_game_year_sequence,
                BatterFullName: batter_full_name_sequence,
                BatterId: batterid_sequence,
                BatterStands: stand_sequence,
                PitcherThrows: throws_sequence,
                Strikes: strike_sequence, 
                Balls: ball_sequence, 
                PitchCount: pitch_count_sequence,
                PitchType: pitch_type_sequence,
                PlateX: plate_x_sequence,
                PlateZ: plate_z_sequence, 
                Zone: zone_sequence,
                Velocity: velocity_sequence,
                IsStrike: is_strike_sequence,
                IsSwing: is_swing_sequence,
                IsContact: is_contact_sequence,
                ExitVelocity: exit_velocity_sequence,
                EventType: event_type_sequence,
          }}

      }

    }

  }  

}
