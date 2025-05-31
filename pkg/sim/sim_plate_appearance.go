package sim

import (
	"fmt"

	"github.com/logananthony/go-baseball/pkg/config"
	"github.com/logananthony/go-baseball/pkg/fetcher"
	"github.com/logananthony/go-baseball/pkg/models"
	"github.com/logananthony/go-baseball/pkg/utils"
)

func SimulatePlateAppearance(pa []models.PlateAppearanceData) []models.PlateAppearanceResult {

	db := config.ConnectDB()
	defer db.Close()

	balls := 0
	strikes := 0
	pitch_count := 0

	leagueSwingProbs := fetcher.FetchBatterSwingPercentageLeague()
	leagueContactProbs := fetcher.FetchBatterContactPercentageLeague()
	leaguePitcherCovMeans := fetcher.FetchPitcherCovarianceMeanLeague()

	batterInfo, err := fetcher.FetchPlayerInfo(db, &pa[0].BatterId, &pa[0].BatterGameYear)
	if err != nil {
	}
	pitcherInfo, err := fetcher.FetchPlayerInfo(db, &pa[0].PitcherId, &pa[0].PitcherGameYear)
	if err != nil {
	}

	if len(batterInfo) == 0 {
		panic(fmt.Sprintf("No batter info found for ID %d", pa[0].BatterId))
	}

	if len(pitcherInfo) == 0 {
		panic(fmt.Sprintf("No pitcher info found for ID %d", pa[0].PitcherId))
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

	batterSwingProbs, _ := fetcher.FetchBatterSwingPercentage(db, pa[0].BatterId, pa[0].BatterGameYear)
	batterContactProbs, _ := fetcher.FetchBatterContactPercentage(db, pa[0].BatterId, pa[0].BatterGameYear)

	batterHitProbs, err := fetcher.FetchBatterHitType(db, pa[0].BatterId, pa[0].BatterGameYear)
	if err != nil {
		fmt.Println("Fetcher error:", err)
	}

	pitcher_game_year_sequence := []int{}
	batter_game_year_sequence := []int{}
	pitcher_full_name_sequence := []string{}
	batter_full_name_sequence := []string{}
	batterid_sequence := []int{}
	pitcherid_sequence := []int{}
	stand_sequence := []string{}
	throws_sequence := []string{}

	pitch_type_sequence := []string{}
	plate_x_sequence := []float64{}
	plate_z_sequence := []float64{}
	zone_sequence := []int{}
	velocity_sequence := []float64{}
	strike_sequence := []int{}
	ball_sequence := []int{}
	pitch_count_sequence := []int{}
	is_strike_sequence := []bool{}
	is_swing_sequence := []bool{}
	is_contact_sequence := []string{}
	exit_velocity_sequence := []float64{}
	launch_angle_sequence := []float64{}
	spray_angle_sequence := []float64{}
	event_type_sequence := []string{}

	away_score_sequence := []int{}
	home_score_sequence := []int{}
	at_bat_number_sequence := []int{}
	inning_sequence := []int{}
	inningtopbot_sequence := []string{}
	outs_sequence := []int{}
	on1b_sequence := []bool{}
	on2b_sequence := []bool{}
	on3b_sequence := []bool{}

	for {

		pitch_count += 1

		pitcher_freqs := fetcher.FetchPitcherFrequencies(db,
			pa[0].PitcherId,
			*batterStands.BatSide)

		pitch_type_result := SimulatePitchType(pitcher_freqs,
			balls,
			strikes)

		pitch_covariance := fetcher.FetchPitcherCovarianceMean(db,
			int64(pa[0].PitcherId),
			int64(pa[0].PitcherGameYear))

		location_velo_result := SimulatePitchLocationVelo(pitch_covariance,
			leaguePitcherCovMeans,
			pitch_type_result,
			*batterStands.BatSide,
			balls,
			strikes)

		zone_result := utils.GetPitchZone(location_velo_result[0],
			location_velo_result[1])

		is_strike_result := utils.IsPitchStrike(location_velo_result[0],
			location_velo_result[1])

		is_swing_result := SimulateSwingDecision(batterSwingProbs,
			leagueSwingProbs,
			*batterStands.BatSide,
			*pitcherThrows.PitchHand,
			pitch_type_result,
			location_velo_result[0],
			location_velo_result[1])

		is_contact_result := SimulateContactPercentage(batterContactProbs,
			leagueContactProbs,
			*batterStands.BatSide,
			*pitcherThrows.PitchHand,
			pitch_type_result,
			location_velo_result[0],
			location_velo_result[1])

		event_type_result := SimulateBatterHitType(batterHitProbs,
			*batterStands.BatSide,
			*pitcherThrows.PitchHand,
			pitch_type_result,
			location_velo_result[0],
			location_velo_result[1],
			location_velo_result[2])

		evRows := fetcher.FetchEVDistributions(
			db,
			pa[0].BatterGameYear,
			pa[0].BatterId,
			*batterStands.BatSide,
			*pitcherThrows.PitchHand,
			utils.StrToNull(event_type_result),
			utils.StrToNull(pitch_type_result),
			utils.IntToNull(zone_result),
			utils.StrToNull(utils.GetVelocityBucket(location_velo_result[2])),
		)
		agg_ev := AggregateEVDistributions(evRows)
		ev_result := SampleFromAggregatedDistribution(agg_ev)

		laRows := fetcher.FetchLADistributions(
			db,
			pa[0].BatterGameYear,
			pa[0].BatterId,
			*batterStands.BatSide,
			*pitcherThrows.PitchHand,
			utils.StrToNull(event_type_result),
			utils.IntToNull(zone_result),
			utils.StrToNull(utils.GetEVBucket(ev_result)),
		)
		agg_la := AggregateLADistributions(laRows)
		la_result := SampleFromAggregatedLADistribution(agg_la)

		sprayRows := fetcher.FetchSprayDistributions(
			db,
			pa[0].BatterGameYear,
			pa[0].BatterId,
			*batterStands.BatSide,
			*pitcherThrows.PitchHand,
			utils.StrToNull(event_type_result),
			utils.IntToNull(zone_result),
			utils.StrToNull(utils.GetEVBucket(ev_result)),
			utils.StrToNull(utils.GetLaunchAngleBucket(la_result)),
		)
		agg_spray := AggregateSprayDistributions(sprayRows)
		spray_result := SampleFromAggregatedSprayDistribution(agg_spray)

		pitcher_game_year_sequence = append(pitcher_game_year_sequence, pa[0].PitcherGameYear)
		batter_game_year_sequence = append(batter_game_year_sequence, pa[0].BatterGameYear)
		pitcher_full_name_sequence = append(pitcher_full_name_sequence, *pitcherThrows.FullName)
		batter_full_name_sequence = append(batter_full_name_sequence, *batterStands.FullName)
		batterid_sequence = append(batterid_sequence, pa[0].BatterId)
		pitcherid_sequence = append(pitcherid_sequence, pa[0].PitcherId)
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
		if !is_swing_result {
			is_contact_sequence = append(is_contact_sequence, "")
		}
		exit_velocity_sequence = append(exit_velocity_sequence, ev_result)
		launch_angle_sequence = append(launch_angle_sequence, la_result)
		spray_angle_sequence = append(spray_angle_sequence, spray_result)
		event_type_sequence = append(event_type_sequence, "")
		away_score_sequence = append(away_score_sequence, pa[0].AwayScore)
		home_score_sequence = append(home_score_sequence, pa[0].HomeScore)
		at_bat_number_sequence = append(at_bat_number_sequence, pa[0].AtBatNumber)
		inning_sequence = append(inning_sequence, pa[0].Inning)
		inningtopbot_sequence = append(inningtopbot_sequence, pa[0].InningTopBot)
		outs_sequence = append(outs_sequence, pa[0].Outs)
		on1b_sequence = append(on1b_sequence, pa[0].On1b)
		on2b_sequence = append(on2b_sequence, pa[0].On2b)
		on3b_sequence = append(on3b_sequence, pa[0].On3b)

		if is_swing_result {
			// Batter swung
			is_contact_sequence = append(is_contact_sequence, is_contact_result)

			switch is_contact_result {
			case "swinging_strike":
				strikes += 1
			case "foul":
				if strikes < 2 {
					strikes += 1
				}
			case "ball_in_play":

				event_type_sequence = append(event_type_sequence, event_type_result)
				event_type_sequence = event_type_sequence[1:]

				return []models.PlateAppearanceResult{{
					PitcherGameYear: pitcher_game_year_sequence,
					PitcherFullName: pitcher_full_name_sequence,
					PitcherId:       pitcherid_sequence,
					BatterGameYear:  batter_game_year_sequence,
					BatterFullName:  batter_full_name_sequence,
					BatterId:        batterid_sequence,
					BatterStands:    stand_sequence,
					PitcherThrows:   throws_sequence,
					Strikes:         strike_sequence,
					Balls:           ball_sequence,
					PitchCount:      pitch_count_sequence,
					PitchType:       pitch_type_sequence,
					PlateX:          plate_x_sequence,
					PlateZ:          plate_z_sequence,
					Zone:            zone_sequence,
					Velocity:        velocity_sequence,
					IsStrike:        is_strike_sequence,
					IsSwing:         is_swing_sequence,
					IsContact:       is_contact_sequence,
					ExitVelocity:    exit_velocity_sequence,
					LaunchAngle:     launch_angle_sequence,
					SprayAngle:      spray_angle_sequence,
					EventType:       event_type_sequence,
					AwayScore:       away_score_sequence,
					HomeScore:       home_score_sequence,
					AtBatNumber:     at_bat_number_sequence,
					Inning:          inning_sequence,
					InningTopBot:    inningtopbot_sequence,
					Outs:            outs_sequence,
					On1b:            on1b_sequence,
					On2b:            on2b_sequence,
					On3b:            on3b_sequence,
				}}

			}

		} else {
			if is_strike_result {
				strikes += 1
			} else {
				balls += 1
			}

		}

		if strikes == 3 || balls == 4 {

			switch {

			case strikes == 3:

				event_type_sequence = append(event_type_sequence, "strikeout")
				event_type_sequence = event_type_sequence[1:]

				return []models.PlateAppearanceResult{{
					PitcherGameYear: pitcher_game_year_sequence,
					PitcherFullName: pitcher_full_name_sequence,
					PitcherId:       pitcherid_sequence,
					BatterGameYear:  batter_game_year_sequence,
					BatterFullName:  batter_full_name_sequence,
					BatterId:        batterid_sequence,
					BatterStands:    stand_sequence,
					PitcherThrows:   throws_sequence,
					Strikes:         strike_sequence,
					Balls:           ball_sequence,
					PitchCount:      pitch_count_sequence,
					PitchType:       pitch_type_sequence,
					PlateX:          plate_x_sequence,
					PlateZ:          plate_z_sequence,
					Zone:            zone_sequence,
					Velocity:        velocity_sequence,
					IsStrike:        is_strike_sequence,
					IsSwing:         is_swing_sequence,
					IsContact:       is_contact_sequence,
					ExitVelocity:    exit_velocity_sequence,
					LaunchAngle:     launch_angle_sequence,
					SprayAngle:      spray_angle_sequence,
					EventType:       event_type_sequence,
					AwayScore:       away_score_sequence,
					HomeScore:       home_score_sequence,
					AtBatNumber:     at_bat_number_sequence,
					Inning:          inning_sequence,
					InningTopBot:    inningtopbot_sequence,
					Outs:            outs_sequence,
					On1b:            on1b_sequence,
					On2b:            on2b_sequence,
					On3b:            on3b_sequence,
				}}
			case balls == 4:

				event_type_sequence = append(event_type_sequence, "walk")
				event_type_sequence = event_type_sequence[1:]

				return []models.PlateAppearanceResult{{
					PitcherGameYear: pitcher_game_year_sequence,
					PitcherFullName: pitcher_full_name_sequence,
					PitcherId:       pitcherid_sequence,
					BatterGameYear:  batter_game_year_sequence,
					BatterFullName:  batter_full_name_sequence,
					BatterId:        batterid_sequence,
					BatterStands:    stand_sequence,
					PitcherThrows:   throws_sequence,
					Strikes:         strike_sequence,
					Balls:           ball_sequence,
					PitchCount:      pitch_count_sequence,
					PitchType:       pitch_type_sequence,
					PlateX:          plate_x_sequence,
					PlateZ:          plate_z_sequence,
					Zone:            zone_sequence,
					Velocity:        velocity_sequence,
					IsStrike:        is_strike_sequence,
					IsSwing:         is_swing_sequence,
					IsContact:       is_contact_sequence,
					ExitVelocity:    exit_velocity_sequence,
					LaunchAngle:     launch_angle_sequence,
					SprayAngle:      spray_angle_sequence,
					EventType:       event_type_sequence,
					AwayScore:       away_score_sequence,
					HomeScore:       home_score_sequence,
					AtBatNumber:     at_bat_number_sequence,
					Inning:          inning_sequence,
					InningTopBot:    inningtopbot_sequence,
					Outs:            outs_sequence,
					On1b:            on1b_sequence,
					On2b:            on2b_sequence,
					On3b:            on3b_sequence,
				}}

			}

		}

	}

}

// func postGameResults(paResult []models.PlateAppearanceResult, db *sql.DB) {
// 	gameId := uuid.New().String()
// 	gameYear := 2024

// 	for _, result := range paResult {

// 		err := poster.InsertGameResult(db, gameId, gameYear, result)
// 		if err != nil {
// 			fmt.Println("Error inserting game result:", err)
// 		}
// 	}
// }
