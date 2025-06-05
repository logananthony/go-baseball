package sim

import (
	"fmt"

	"github.com/logananthony/go-baseball/pkg/config"
	"github.com/logananthony/go-baseball/pkg/models"
	"github.com/logananthony/go-baseball/pkg/utils"
)

func SimulatePlateAppearance(pa []models.PlateAppearanceData, sim []models.SimData) []models.PlateAppearanceResult {

	db := config.ConnectDB()
	defer db.Close()

	playerInfoMap := make(map[int]models.MLBPlayerInfo)
	for _, player := range sim[0].PlayerInfo {
		if player.ID != nil {
			playerInfoMap[*player.ID] = player
		}
	}

	batterInfo, ok := playerInfoMap[pa[0].BatterId]
	if !ok {
		panic(fmt.Sprintf("Missing batter info for ID: %d", pa[0].BatterId))
	}

	pitcherInfo, ok := playerInfoMap[pa[0].PitcherId]
	if !ok {
		panic(fmt.Sprintf("Missing pitcher info for ID: %d", pa[0].PitcherId))
	}

	batterStands := batterInfo.BatSide
	pitcherThrows := pitcherInfo.PitchHand

	swingPctMap := make(map[int]models.BatterSwingPercentage)
	swingPctStore := []models.BatterSwingPercentage{}
	for _, player := range sim[0].BatterSwing {
		swingPctMap[player.Batter] = player
		swingPctStore = append(swingPctStore, player)
	}

	contactPctMap := make(map[int]models.BatterContactPercentage)
	contactPctStore := []models.BatterContactPercentage{}
	for _, player := range sim[0].BatterContact {
		contactPctMap[player.Batter] = player
		contactPctStore = append(contactPctStore, player)
	}

	hitProbsMap := make(map[int]models.BatterHitType)
	hitProbsStore := []models.BatterHitType{}
	for _, player := range sim[0].BatterHitType {
		hitProbsMap[player.Batter] = player
		hitProbsStore = append(hitProbsStore, player)
	}

	balls := 0
	strikes := 0
	pitch_count := 0

	if batterStands != nil {
		switch *batterStands {
		case "Left":
			left := "L"
			batterStands = &left
		case "Right":
			right := "R"
			batterStands = &right
		case "Switch":
			if pitcherThrows != nil {
				if *pitcherThrows == "Right" {
					left := "L"
					batterStands = &left
				} else if *pitcherThrows == "Left" {
					right := "R"
					batterStands = &right
				}
			}
		}
	}

	// fmt.Println("Batter Stands:", *batterStands)

	if pitcherThrows != nil {
		switch *pitcherThrows {
		case "Left":
			left := "L"
			pitcherThrows = &left
		case "Right":
			right := "R"
			pitcherThrows = &right
		}
	}

	// batterSwingProbs, _ := fetcher.FetchBatterSwingPercentage(db, pa[0].BatterId, pa[0].BatterGameYear)
	// batterContactProbs, _ := fetcher.FetchBatterContactPercentage(db, pa[0].BatterId, pa[0].BatterGameYear)
	// batterHitProbs, err := fetcher.FetchBatterHitType(db, pa[0].BatterId, pa[0].BatterGameYear)
	// if err != nil {
	// 	fmt.Println("Fetcher error:", err)
	// }

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

		var pitcher_freqs []models.PitcherCountPitchFreq
		for _, player := range sim[0].PitcherPitchFreq {
			if player.PITCHER == pa[0].PitcherId && player.STAND == *batterStands {
				pitcher_freqs = append(pitcher_freqs, player)
			}
		}

		pitch_type_result := SimulatePitchType(pitcher_freqs, balls, strikes)

		// fmt.Println("SimulatePitchType result: ", pitch_type_result)

		pitcherCovMap := make(map[int]models.PitcherCovarianceMean)
		var pitcherCovMeanStore []models.PitcherCovarianceMean
		for _, player := range sim[0].PitcherCovMeans {
			pitcherCovMap[int(player.Pitcher)] = player
			pitcherCovMeanStore = append(pitcherCovMeanStore, player)
		}

		location_velo_result := SimulatePitchLocationVelo(pitcherCovMeanStore,
			sim[0].LeaguePitchCovMeans,
			pitch_type_result,
			*batterStands,
			balls,
			strikes)

		// fmt.Println("Location Velo Result:", location_velo_result)

		// fmt.Println("Pitch Covariance:", pitch_covariance)
		// fmt.Println("League Pitch Cov Means:", sim[0].LeaguePitchCovMeans)
		// fmt.Println("Pitch Type Result:", pitch_type_result)
		// fmt.Println("Batter Stands:", *batterStands)
		// fmt.Println("Balls:", balls, "Strikes:", strikes)

		zone_result := utils.GetPitchZone(location_velo_result[0],
			location_velo_result[1])

		// fmt.Println(zone_result)

		is_strike_result := utils.IsPitchStrike(location_velo_result[0],
			location_velo_result[1])

		is_swing_result := SimulateSwingDecision(swingPctStore,
			sim[0].LeagueSwing,
			*batterStands,
			*pitcherThrows,
			pitch_type_result,
			location_velo_result[0],
			location_velo_result[1])

		// fmt.Println("Is Swing Result:", is_swing_result)

		is_contact_result := SimulateContactPercentage(contactPctStore,
			sim[0].LeagueContact,
			*batterStands,
			*pitcherThrows,
			pitch_type_result,
			location_velo_result[0],
			location_velo_result[1])

		event_type_result := SimulateBatterHitType(hitProbsStore,
			*batterStands,
			*pitcherThrows,
			pitch_type_result,
			location_velo_result[0],
			location_velo_result[1],
			location_velo_result[2])

		// fmt.Println("Event Type Result:", event_type_result)
		// fmt.println(strikes)

		batterEvMap := make(map[int]models.EVDistribution)
		batterEVStore := []models.EVDistribution{}
		for _, player := range sim[0].BatterEVDist {
			batterEvMap[int(player.Batter)] = player
			batterEVStore = append(batterEVStore, player)
		}
		var evRows models.EVDistribution
		for _, player := range batterEVStore {
			if player.Batter == pa[0].BatterId &&
				player.GameYear == pa[0].BatterGameYear &&
				player.Stand == *batterStands &&
				player.PThrows == *pitcherThrows &&
				player.Outcome == utils.StrToNull(event_type_result) &&
				player.PitchType == utils.StrToNull(pitch_type_result) &&
				player.Zone == utils.IntToNull(zone_result) &&
				player.VelocityBucket == utils.StrToNull(utils.GetVelocityBucket(location_velo_result[2])) {
				evRows = player
				break
			}
		}

		agg_ev := AggregateEVDistributions([]models.EVDistribution{evRows})
		ev_result := SampleFromAggregatedDistribution(agg_ev)

		// fmt.Println("EV Result:", ev_result)

		laDistMap := make(map[int]models.LADistribution)
		laDistStore := []models.LADistribution{}
		for _, player := range sim[0].BatterLADist {
			laDistMap[int(player.Batter)] = player
			laDistStore = append(laDistStore, player)
		}

		var laRows models.LADistribution
		for _, player := range laDistStore {
			if player.Batter == pa[0].BatterId &&
				player.GameYear == pa[0].BatterGameYear &&
				player.Stand == *batterStands &&
				player.PThrows == *pitcherThrows &&
				player.Outcome == &event_type_result &&
				player.Zone == &zone_result {

				evBucket := utils.GetEVBucket(ev_result)
				if player.EVBucket != nil && *player.EVBucket == evBucket {
					laRows = player
					break
				}
			}
		}

		agg_la := AggregateLADistributions([]models.LADistribution{laRows})
		la_result := SampleFromAggregatedLADistribution(agg_la)

		sprayDistMap := make(map[int]models.SprayDistribution)
		sprayDistStore := []models.SprayDistribution{}
		for _, player := range sim[0].BatterSprayDist {
			sprayDistMap[int(player.Batter)] = player
			sprayDistStore = append(sprayDistStore, player)
		}

		var sprayRows models.SprayDistribution
		for _, player := range sprayDistStore {
			if player.Batter == pa[0].BatterId &&
				player.GameYear == pa[0].BatterGameYear &&
				player.Stand == *batterStands &&
				player.PThrows == *pitcherThrows &&
				player.Outcome == &event_type_result &&
				player.Zone == &zone_result &&
				player.EVBucket != nil && *player.EVBucket == utils.GetEVBucket(ev_result) &&
				player.LaunchAngleBucket != nil && *player.LaunchAngleBucket == utils.GetLaunchAngleBucket(la_result) {
				sprayRows = player
				break
			}
		}

		agg_spray := AggregateSprayDistributions([]models.SprayDistribution{sprayRows})
		spray_result := SampleFromAggregatedSprayDistribution(agg_spray)

		pitcher_game_year_sequence = append(pitcher_game_year_sequence, pa[0].PitcherGameYear)
		batter_game_year_sequence = append(batter_game_year_sequence, pa[0].BatterGameYear)
		pitcher_full_name_sequence = append(pitcher_full_name_sequence, *pitcherInfo.FullName)
		batter_full_name_sequence = append(batter_full_name_sequence, *batterInfo.FullName)
		batterid_sequence = append(batterid_sequence, pa[0].BatterId)
		pitcherid_sequence = append(pitcherid_sequence, pa[0].PitcherId)
		stand_sequence = append(stand_sequence, *batterStands)
		throws_sequence = append(throws_sequence, *pitcherThrows)
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

		// log.Printf("Pitch #%d | Type=%s | PlateX=%.2f PlateZ=%.2f | Zone=%d | Velo=%.1f | Strike=%t | Swing=%t",
		// 	pitch_count, pitch_type_result, location_velo_result[0], location_velo_result[1], zone_result, location_velo_result[2],
		// 	is_strike_result, is_swing_result)

		if is_swing_result {
			// Batter swung
			is_contact_sequence = append(is_contact_sequence, is_contact_result)

			// if is_swing_result {
			// 	log.Printf("Swing Decision: ContactType=%s", is_contact_result)
			// 	if is_contact_result == "ball_in_play" {
			// 		log.Printf("Ball In Play: Event=%s | EV=%.1f LA=%.1f SA=%.1f",
			// 			event_type_result, ev_result, la_result, spray_result)
			// 	}
			// }

			// if !is_swing_result {
			// 	if is_strike_result {
			// 		log.Printf("Called Strike")
			// 	} else {
			// 		log.Printf("Ball")
			// 	}
			// }

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
				// fmt.Println("Event Type:", event_type_sequence)

				// log.Printf("PA End: Result=%s | Final Balls=%d, Strikes=%d, Total Pitches=%d",
				// 	event_type_sequence[len(event_type_sequence)-1], balls, strikes, pitch_count)

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

				// log.Printf("PA End: Result=%s | Final Balls=%d, Strikes=%d, Total Pitches=%d",
				// 	event_type_sequence[len(event_type_sequence)-1], balls, strikes, pitch_count)

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
