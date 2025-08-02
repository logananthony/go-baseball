package sim

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/logananthony/go-baseball/pkg/fetcher"
	"github.com/logananthony/go-baseball/pkg/models"
	"github.com/logananthony/go-baseball/pkg/poster"
	"github.com/logananthony/go-baseball/pkg/utils"
)

// func SimulateGame(db *sql.DB, gameData []models.GameData) {
func SimulateGame(
	db *sql.DB,
	simData models.SimData,
	homeLineup []models.GameDataGamePk,
	awayLineup []models.GameDataGamePk,
	pitchingSubProbs []models.PitchingSubstitutionProb,
	homeBullpen *models.BullpenOrder,
	awayBullpen *models.BullpenOrder,
	bullpenRoleProbs []models.BullpenRoleProb,
	gameData models.GameData,
	outcomeMode string,
) (homeScore int, awayScore int, homeTeam string, awayTeam string,
	homeStartingPitcherName string, awayStartingPitcherName string,
	gameDate time.Time, batterEvents []models.BatterEvent, pitcherEvents []models.PitcherEvent, runEvents []models.RunEvent) {

	defer func() {
		if r := recover(); r != nil {
			log.Printf("🔥 Panic recovered in SimulateGame: %v", r)
		}
	}()

	gamePk, _ := fetcher.FetchGameDataByGamePk(db, gameData.GamePk)

	first := gamePk[0]

	gameDate = first.GameDate
	homeStartingPitcher := first.HomePitcherId
	awayStartingPitcher := first.AwayPitcherId
	homeStartingPitcherName = first.HomePitcherName
	awayStartingPitcherName = first.AwayPitcherName
	homeTeam = first.HomeTeamAbbr
	awayTeam = first.AwayTeamAbbr
	season := first.Season

	gameRes := models.GameResult{
		GameId: uuid.New().String(),
		GamePk: gameData.GamePk,
		JobId:  gameData.JobId,
	}

	homePitcherLineup := [][]int{
		{homeStartingPitcher, season},
		{homeBullpen.PlayerID1, season},
		{homeBullpen.PlayerID2, season},
		{homeBullpen.PlayerID3, season},
		{homeBullpen.PlayerID4, season},
		{homeBullpen.PlayerID5, season},
		{homeBullpen.PlayerID6, season},
		{homeBullpen.PlayerID7, season},
		{homeBullpen.PlayerID8, season},
	}

	if awayBullpen == nil {
		log.Fatalf("awayBullpen is nil for team: %s, season: %d", awayTeam, season)
	}
	awayPitcherLineup := [][]int{
		{awayStartingPitcher, season},
		{awayBullpen.PlayerID1, season},
		{awayBullpen.PlayerID2, season},
		{awayBullpen.PlayerID3, season},
		{awayBullpen.PlayerID4, season},
		{awayBullpen.PlayerID5, season},
		{awayBullpen.PlayerID6, season},
		{awayBullpen.PlayerID7, season},
		{awayBullpen.PlayerID8, season},
	}

	var homePitcher int
	var homePitcherGameYear int
	var awayPitcher int
	var awayPitcherGameYear int

	inning := 1
	awayScore = 0
	homeScore = 0
	awayBatterNumber := 0
	homeBatterNumber := 0
	atBatNumber := 0
	homePitcher = homePitcherLineup[0][0]
	homePitcherGameYear = homePitcherLineup[0][1]
	awayPitcher = awayPitcherLineup[0][0]
	awayPitcherGameYear = awayPitcherLineup[0][1]

	// fmt.Println("Caching league-level data...")

	simData.HomeTeam = homeTeam

	homePitchCount := 0
	awayPitchCount := 0
	var pitcherPulledHome bool
	var pitcherPulledAway bool

	// Track per-game appearances and totals
	appearedPitcher := make(map[int64]bool)
	strikeoutsByPitcher := make(map[int64]int)
	swstrByPitcher := make(map[int64]int) // ✅ was unused before; we’ll fill it

	appearedBatter := make(map[int64]bool)
	runsByBatter := make(map[int64]int)
	rbiByBatter := make(map[int64]int)   // ✅ new
	swstrByBatter := make(map[int64]int) // ✅ new

	// Helper to append appearance-based events before returning
	// Helper to append appearance-based events before returning
	finalizeAppearanceEvents := func() {
		// Pitchers: appeared but 0 Ks
		for pid := range appearedPitcher {
			if strikeoutsByPitcher[pid] == 0 {
				pitcherEvents = append(pitcherEvents, models.PitcherEvent{
					PitcherID: pid,
					EventType: "appeared_no_k",
				})
			}
			// ✅ Pitchers: appeared but 0 swinging strikes
			if swstrByPitcher[pid] == 0 {
				pitcherEvents = append(pitcherEvents, models.PitcherEvent{
					PitcherID: pid,
					EventType: "appeared_no_swstr",
				})
			}
		}

		// Batters: appeared but 0 runs (existing)
		for bid := range appearedBatter {
			if runsByBatter[bid] == 0 {
				batterEvents = append(batterEvents, models.BatterEvent{
					BatterID:  bid,
					EventType: "appeared_no_run",
					Runs:      0,
				})
			}

			// ✅ Batters: appeared but 0 RBI
			if rbiByBatter[bid] == 0 {
				batterEvents = append(batterEvents, models.BatterEvent{
					BatterID:  bid,
					EventType: "appeared_no_rbi",
					RBI:       0,
				})
			}

			// ✅ Batters: appeared but 0 swinging strikes
			if swstrByBatter[bid] == 0 {
				batterEvents = append(batterEvents, models.BatterEvent{
					BatterID:   bid,
					EventType:  "appeared_no_swstr",
					SwStrCount: 0,
				})
			}
		}
	}

	for {

		topOuts := 0
		botOuts := 0
		awayBaseState := []int{0, 0, 0}
		homeBaseState := []int{0, 0, 0}
		var scoredRuns []int
		var rbi int
		// priorAwayScore := awayScore
		// priorHomeScore := homeScore

		// inningRunsHome := awayScore - priorAwayScore
		// inningRunsAway := homeScore - priorHomeScore

		for topOuts < 3 {
			awayBatterNumber = awayBatterNumber % 9
			awayBatter := awayLineup[awayBatterNumber]
			awayBatterGameYear := awayBatter.Season
			awayPaResult := SimulatePlateAppearance([]models.PlateAppearanceData{{
				BatterGameYear:  awayBatterGameYear,
				BatterId:        awayBatter.PlayerId,
				PitcherGameYear: homePitcherGameYear,
				PitcherId:       homePitcher,
				Strikes:         0,
				Balls:           0,
				AwayScore:       awayScore,
				HomeScore:       homeScore,
				AtBatNumber:     atBatNumber,
				Inning:          inning,
				InningTopBot:    "Top",
				Outs:            topOuts,
				On1b:            awayBaseState[0] != 0,
				On2b:            awayBaseState[1] != 0,
				On3b:            awayBaseState[2] != 0,
				RunnerOn1bID:    awayBaseState[0],
				RunnerOn2bID:    awayBaseState[1],
				RunnerOn3bID:    awayBaseState[2],
			}}, []models.SimData{simData})

			atBatNumber++

			// fmt.Println("Home Pitch Count:", homePitchCount)
			awayScore, awayBaseState, topOuts, homePitchCount, scoredRuns, rbi = ProcessPlateAppearance(
				awayPaResult, awayScore, awayBaseState, topOuts, homePitchCount,
			)

			// After ProcessPlateAppearance(...)
			for _, id := range scoredRuns {
				// keep: internal per-game tally
				runsByBatter[int64(id)]++
				// add: emit event for API aggregation
				runEvents = append(runEvents, models.RunEvent{
					RunnerID: int64(id), // only RunnerID, unless you add fields to the struct
				})
			}

			// credit total runs for whoever scored in this PA
			for _, id := range scoredRuns {
				runsByBatter[int64(id)]++
			}
			// ✅ mark batter appeared using the last batter id from the PA
			if len(awayPaResult) > 0 && len(awayPaResult[0].BatterId) > 0 {
				lastBatterID := int64(awayPaResult[0].BatterId[len(awayPaResult[0].BatterId)-1])
				appearedBatter[lastBatterID] = true
			}

			for _, paResult := range awayPaResult {

				AppendPlateAppearanceTopResult(paResult, awayScore, homeScore, atBatNumber, inning, topOuts, awayBaseState)
				AppendGameResult(&gameRes, paResult, awayBatter.PlayerId, rbi)

				// mark pitcher appeared (use last pitcher id for this PA)
				if len(paResult.PitcherId) > 0 {
					pid := int64(paResult.PitcherId[len(paResult.PitcherId)-1])
					appearedPitcher[pid] = true
				}

				if len(paResult.BatterId) > 0 && len(paResult.EventType) > 0 {
					lastBatterID := int64(paResult.BatterId[len(paResult.BatterId)-1])
					lastPitcherID := int64(paResult.PitcherId[len(paResult.PitcherId)-1])

					batterRuns := 0
					for _, id := range scoredRuns {
						if int64(id) == lastBatterID {
							batterRuns++
						}
					}

					// ✅ count swinging strikes in this PA
					// thisPASwStr := int64(utils.IsSwingingStrikeCount(paResult.IsSwing, paResult.IsStrike))
					thisPASwStr := "swinging_strike"

					swStrCount := 0
					for _, f := range paResult.IsContact {
						if f == thisPASwStr {
							swStrCount++
						}
					}

					batterEvents = append(batterEvents, models.BatterEvent{
						BatterID:   lastBatterID,
						PitcherID:  lastPitcherID,
						EventType:  paResult.EventType[len(paResult.EventType)-1],
						Runs:       batterRuns,
						RBI:        rbi,
						SwStrCount: swStrCount, // ✅ include per-PA value
						PitchCount: len(paResult.BatterId),
					})

					// ✅ accumulate per-game tallies
					rbiByBatter[lastBatterID] += rbi
					swstrByBatter[lastBatterID] += swStrCount

					// ✅ also attribute swinging strikes to the pitcher
					swstrByPitcher[lastPitcherID] += swStrCount
				}

				// strikeout event for pitchers (use last pitcher id)
				if len(paResult.EventType) > 0 && paResult.EventType[len(paResult.EventType)-1] == "strikeout" && len(paResult.PitcherId) > 0 {
					pid := int64(paResult.PitcherId[len(paResult.PitcherId)-1])
					strikeoutsByPitcher[pid]++
					pitcherEvents = append(pitcherEvents, models.PitcherEvent{PitcherID: pid, EventType: "strikeout"})
				}

			}

			awayBatterNumber++

		}

		if inning >= 9 && homeScore > awayScore {
			if outcomeMode == "pbp" {
				postGameResults([]models.GameResult{gameRes}, season, db)
			}
			finalizeAppearanceEvents()
			fmt.Println("Home team wins:", homeScore, "-", awayScore)
			return homeScore, awayScore, homeTeam, awayTeam, homeStartingPitcherName, awayStartingPitcherName, gameDate, batterEvents, pitcherEvents, runEvents
		}

		usedPitchersHome := map[int]bool{}

		if pitcherPulledHome {
			if len(homePitcherLineup) > 0 {
				homePitcherLineup = utils.FilterSliceSlices(homePitcherLineup, homePitcher)
				if len(homePitcherLineup) > 0 {
					runDiff := homeScore - awayScore
					runnersOn := utils.CountRunners(awayBaseState)

					selected := utils.SelectBullpenPitcherLineup(
						homePitcherLineup,
						bullpenRoleProbs,
						inning,
						runDiff,
						runnersOn,
						usedPitchersHome,
					)
					// fmt.Println("pitcherPulledHome:", pitcherPulledHome)
					// fmt.Println("homePitcherLineup:", homePitcherLineup)
					// fmt.Println("bullpenRoleProbs:", bullpenRoleProbs)
					// fmt.Println("inning:", inning)
					// fmt.Println("runDiff:", runDiff)
					// fmt.Println("runnersOn:", runnersOn)
					// fmt.Println("usedPitchersHome:", usedPitchersHome)
					// fmt.Println("Selected Pitcher: ", selected)

					if selected != nil {
						homePitcher = selected[0]
						homePitcherGameYear = selected[1]
						usedPitchersHome[homePitcher] = true
						// fmt.Printf("Home Pitch Count Before: %d", homePitchCount)
						homePitchCount = 0
						// fmt.Printf("Home Pitch Count After: %d", homePitchCount)
					} else {
						// fmt.Println("No eligible pitchers left for inning", inning)
					}

				} else {
					// fmt.Println("Home pitcher lineup is empty after filtering, skipping pitcher substitution.")
				}
			} else {
				// fmt.Println("Home pitcher lineup is empty, skipping pitcher substitution.")
			}
		}
		homeTeamDeficit := homeScore - awayScore
		// fmt.Println("Home Team Deficit: %d", homeTeamDeficit)

		homePitchBin := utils.PitchCountToBin(homePitchCount)
		// fmt.Println("Home Pitch Count:", homePitchCount)
		// fmt.Println("Home Pitch Bin:", homePitchBin)

		var pitcherPulledHomeProb float64

		if homePitcher == homeStartingPitcher {
			pitcherPulledHomeProb, _ = utils.FilterPitcherExitDistributions(
				pitchingSubProbs,
				"starter",
				homeTeamDeficit,
				homePitchBin,
			)
			// fmt.Printf("pitchingSubProbs: %+v\n", pitchingSubProbs)
			// fmt.Printf("role: %s\n", "starter")
			// fmt.Printf("homeTeamDeficit: %v\n", homeTeamDeficit)
			// fmt.Printf("homePitchBin: %v\n", homePitchBin)
			// fmt.Printf("pitcherPulledHomeProb: %v\n", pitcherPulledHomeProb)
		} else {
			pitcherPulledHomeProb, _ = utils.FilterPitcherExitDistributions(
				pitchingSubProbs,
				"relief",
				homeTeamDeficit,
				homePitchBin,
			)
		}

		if pitcherPulledHomeProb != 0.0 {

			// fmt.Printf("pitcherPulledHomeProb: %v\n", &pitcherPulledHomeProb)
			pitcherPulledHome = utils.IsSuccess(&pitcherPulledHomeProb)
			// fmt.Println("pitcherPulledHome: ", pitcherPulledHome)
		} else {
			pitcherPulledHome = false
		}

		for botOuts < 3 {
			homeBatterNumber = homeBatterNumber % 9
			homeBatter := homeLineup[homeBatterNumber]
			homeBatterGameYear := homeBatter.Season
			homePaResult := SimulatePlateAppearance([]models.PlateAppearanceData{{
				BatterGameYear:  homeBatterGameYear,
				BatterId:        homeBatter.PlayerId,
				PitcherGameYear: awayPitcherGameYear,
				PitcherId:       awayPitcher,
				Strikes:         0,
				Balls:           0,
				AwayScore:       awayScore,
				HomeScore:       homeScore,
				AtBatNumber:     atBatNumber,
				Inning:          inning,
				InningTopBot:    "Bot",
				Outs:            topOuts,
				On1b:            homeBaseState[0] != 0,
				On2b:            homeBaseState[1] != 0,
				On3b:            homeBaseState[2] != 0,
				RunnerOn1bID:    homeBaseState[0],
				RunnerOn2bID:    homeBaseState[1],
				RunnerOn3bID:    homeBaseState[2],
			}}, []models.SimData{simData})

			atBatNumber++

			homeScore, homeBaseState, botOuts, awayPitchCount, scoredRuns, rbi = ProcessPlateAppearance(
				homePaResult, homeScore, homeBaseState, botOuts, awayPitchCount,
			)

			// After ProcessPlateAppearance(...)
			for _, id := range scoredRuns {
				// keep: internal per-game tally
				runsByBatter[int64(id)]++
				// add: emit event for API aggregation
				runEvents = append(runEvents, models.RunEvent{
					RunnerID: int64(id), // only RunnerID, unless you add fields to the struct
				})
			}

			for _, id := range scoredRuns {
				runsByBatter[int64(id)]++
			}
			// mark batter appeared
			if len(homePaResult) > 0 && len(homePaResult[0].BatterId) > 0 {
				appearedBatter[int64(homePaResult[0].BatterId[len(homePaResult[0].BatterId)-1])] = true
			}

			for _, paResult := range homePaResult {
				AppendPlateAppearanceBotResult(paResult, awayScore, homeScore, atBatNumber, inning, botOuts, homeBaseState)
				AppendGameResult(&gameRes, paResult, homeBatter.PlayerId, rbi)

				// mark pitcher appeared (use last pitcher id for this PA)
				if len(paResult.PitcherId) > 0 {
					pid := int64(paResult.PitcherId[len(paResult.PitcherId)-1])
					appearedPitcher[pid] = true
				}

				if len(paResult.BatterId) > 0 && len(paResult.EventType) > 0 {
					lastBatterID := int64(paResult.BatterId[len(paResult.BatterId)-1])
					lastPitcherID := int64(paResult.PitcherId[len(paResult.PitcherId)-1])

					batterRuns := 0
					for _, id := range scoredRuns {
						if int64(id) == lastBatterID {
							batterRuns++
						}
					}

					thisPASwStr := "swinging_strike"

					swStrCount := 0
					for _, f := range paResult.IsContact {
						if f == thisPASwStr {
							swStrCount++
						}
					}

					batterEvents = append(batterEvents, models.BatterEvent{
						BatterID:   lastBatterID,
						PitcherID:  lastPitcherID,
						EventType:  paResult.EventType[len(paResult.EventType)-1],
						Runs:       batterRuns,
						RBI:        rbi,
						SwStrCount: swStrCount, // ✅ include per-PA value
						PitchCount: len(paResult.BatterId),
					})

					// ✅ accumulate per-game tallies
					rbiByBatter[lastBatterID] += rbi
					swstrByBatter[lastBatterID] += swStrCount

					// ✅ also attribute swinging strikes to the pitcher
					swstrByPitcher[lastPitcherID] += swStrCount
				}

				// strikeout event for pitchers (use last pitcher id)
				if len(paResult.EventType) > 0 && paResult.EventType[len(paResult.EventType)-1] == "strikeout" && len(paResult.PitcherId) > 0 {
					pid := int64(paResult.PitcherId[len(paResult.PitcherId)-1])
					strikeoutsByPitcher[pid]++
					pitcherEvents = append(pitcherEvents, models.PitcherEvent{PitcherID: pid, EventType: "strikeout"})
				}

			}

			// spew.Dump(gameRes.PAResult)
			homeBatterNumber++

			if inning >= 9 && homeScore > awayScore {
				fmt.Println("Home team wins (walk-off):", homeScore, "-", awayScore)
				finalizeAppearanceEvents()
				if outcomeMode == "pbp" {
					postGameResults([]models.GameResult{gameRes}, season, db)
				}
				return homeScore, awayScore, homeTeam, awayTeam, homeStartingPitcherName, awayStartingPitcherName, gameDate, batterEvents, pitcherEvents, runEvents
			}
		}

		usedPitchersAway := map[int]bool{}

		if pitcherPulledAway {
			if len(awayPitcherLineup) > 0 {
				awayPitcherLineup = utils.FilterSliceSlices(awayPitcherLineup, awayPitcher)
				if len(awayPitcherLineup) > 0 {
					runDiff := awayScore - homeScore
					runnersOn := utils.CountRunners(homeBaseState)

					selected := utils.SelectBullpenPitcherLineup(
						awayPitcherLineup,
						bullpenRoleProbs,
						inning,
						runDiff,
						runnersOn,
						usedPitchersAway,
					)
					// fmt.Println("pitcherPulledAway:", pitcherPulledAway)
					// fmt.Println("awayPitcherLineup:", awayPitcherLineup)
					// fmt.Println("bullpenRoleProbs:", bullpenRoleProbs)
					// fmt.Println("inning:", inning)
					// fmt.Println("runDiff:", runDiff)
					// fmt.Println("runnersOn:", runnersOn)
					// fmt.Println("usedPitchersAway:", usedPitchersAway)
					// fmt.Println("Selected Pitcher: ", selected)

					if selected != nil {
						awayPitcher = selected[0]
						awayPitcherGameYear = selected[1]
						usedPitchersAway[awayPitcher] = true
						// fmt.Printf("Away Pitch Count Before: %d", awayPitchCount)
						awayPitchCount = 0
						// fmt.Printf("Away Pitch Count After: %d", awayPitchCount)
					} else {
						// fmt.Println("No eligible pitchers left for inning", inning)
					}

				} else {
					// fmt.Println("Away pitcher lineup is empty after filtering, skipping pitcher substitution.")
				}
			} else {
				// fmt.Println("Away pitcher lineup is empty, skipping pitcher substitution.")
			}
		}
		awayTeamDeficit := awayScore - homeScore
		// fmt.Println("Away Team Deficit: %d", awayTeamDeficit)

		awayPitchBin := utils.PitchCountToBin(awayPitchCount)
		// fmt.Println("Away Pitch Count:", awayPitchCount)
		// fmt.Println("Away Pitch Bin:", awayPitchBin)

		var pitcherPulledAwayProb float64

		if awayPitcher == awayStartingPitcher {
			pitcherPulledAwayProb, _ = utils.FilterPitcherExitDistributions(
				pitchingSubProbs,
				"starter",
				awayTeamDeficit,
				awayPitchBin,
			)
			// fmt.Printf("pitchingSubProbs: %+v\n", pitchingSubProbs)
			// fmt.Printf("role: %s\n", "starter")
			// fmt.Printf("awayTeamDeficit: %v\n", awayTeamDeficit)
			// fmt.Printf("awayPitchBin: %v\n", awayPitchBin)
			// fmt.Printf("pitcherPulledAwayProb: %v\n", pitcherPulledAwayProb)
		} else {
			pitcherPulledAwayProb, _ = utils.FilterPitcherExitDistributions(
				pitchingSubProbs,
				"relief",
				awayTeamDeficit,
				awayPitchBin,
			)
		}

		if pitcherPulledAwayProb != 0.0 {

			// fmt.Printf("pitcherPulledAwayProb: %v\n", &pitcherPulledAwayProb)
			pitcherPulledAway = utils.IsSuccess(&pitcherPulledAwayProb)
			// fmt.Println("pitcherPulledAway: ", pitcherPulledAway)
		} else {
			pitcherPulledAway = false
		}

		// If 9 or later and not tied, game ends
		if inning >= 9 && homeScore != awayScore {
			fmt.Println("Away team wins:", awayScore, "-", homeScore)
			finalizeAppearanceEvents()
			if outcomeMode == "pbp" {
				postGameResults([]models.GameResult{gameRes}, season, db)
			}
			return homeScore, awayScore, homeTeam, awayTeam, homeStartingPitcherName, awayStartingPitcherName, gameDate, batterEvents, pitcherEvents, runEvents
		}

		inning++
	}

}

func AppendPlateAppearanceTopResult(paResult models.PlateAppearanceResult, awayScore int, homeScore int, atBatNumber int, inning int, topOuts int, awayBaseState []int) {
	paResult.AwayScore = append(paResult.AwayScore, awayScore)
	paResult.HomeScore = append(paResult.HomeScore, homeScore)
	paResult.AtBatNumber = append(paResult.AtBatNumber, atBatNumber)
	paResult.Inning = append(paResult.Inning, inning)
	paResult.InningTopBot = append(paResult.InningTopBot, "Top")
	paResult.Outs = append(paResult.Outs, topOuts)
	paResult.On1b = append(paResult.On1b, awayBaseState[0] != 0)
	paResult.On2b = append(paResult.On2b, awayBaseState[1] != 0)
	paResult.On3b = append(paResult.On3b, awayBaseState[2] != 0)
}

func AppendPlateAppearanceBotResult(paResult models.PlateAppearanceResult, awayScore int, homeScore int, atBatNumber int, inning int, topOuts int, homeBaseState []int) {
	paResult.AwayScore = append(paResult.AwayScore, awayScore)
	paResult.HomeScore = append(paResult.HomeScore, homeScore)
	paResult.AtBatNumber = append(paResult.AtBatNumber, atBatNumber)
	paResult.Inning = append(paResult.Inning, inning)
	paResult.InningTopBot = append(paResult.InningTopBot, "Bot")
	paResult.Outs = append(paResult.Outs, topOuts)
	paResult.On1b = append(paResult.On1b, homeBaseState[0] != 0)
	paResult.On2b = append(paResult.On2b, homeBaseState[1] != 0)
	paResult.On3b = append(paResult.On3b, homeBaseState[2] != 0)
}

func AppendGameResult(gameRes *models.GameResult, paResult models.PlateAppearanceResult, rbiBatterID int, rbiCount int) {
	gameRes.PAResult.PitcherGameYear = append(gameRes.PAResult.PitcherGameYear, paResult.PitcherGameYear...)
	gameRes.PAResult.PitcherFullName = append(gameRes.PAResult.PitcherFullName, paResult.PitcherFullName...)
	gameRes.PAResult.PitcherId = append(gameRes.PAResult.PitcherId, paResult.PitcherId...)
	gameRes.PAResult.BatterGameYear = append(gameRes.PAResult.BatterGameYear, paResult.BatterGameYear...)
	gameRes.PAResult.BatterFullName = append(gameRes.PAResult.BatterFullName, paResult.BatterFullName...)
	gameRes.PAResult.BatterId = append(gameRes.PAResult.BatterId, paResult.BatterId...)
	gameRes.PAResult.BatterStands = append(gameRes.PAResult.BatterStands, paResult.BatterStands...)
	gameRes.PAResult.PitcherThrows = append(gameRes.PAResult.PitcherThrows, paResult.PitcherThrows...)
	gameRes.PAResult.Strikes = append(gameRes.PAResult.Strikes, paResult.Strikes...)
	gameRes.PAResult.Balls = append(gameRes.PAResult.Balls, paResult.Balls...)
	gameRes.PAResult.PitchCount = append(gameRes.PAResult.PitchCount, paResult.PitchCount...)
	gameRes.PAResult.PitchType = append(gameRes.PAResult.PitchType, paResult.PitchType...)
	gameRes.PAResult.PlateX = append(gameRes.PAResult.PlateX, paResult.PlateX...)
	gameRes.PAResult.PlateZ = append(gameRes.PAResult.PlateZ, paResult.PlateZ...)
	gameRes.PAResult.Zone = append(gameRes.PAResult.Zone, paResult.Zone...)
	gameRes.PAResult.Velocity = append(gameRes.PAResult.Velocity, paResult.Velocity...)
	gameRes.PAResult.IsStrike = append(gameRes.PAResult.IsStrike, paResult.IsStrike...)
	gameRes.PAResult.IsSwing = append(gameRes.PAResult.IsSwing, paResult.IsSwing...)
	gameRes.PAResult.IsContact = append(gameRes.PAResult.IsContact, paResult.IsContact...)
	gameRes.PAResult.ExitVelocity = append(gameRes.PAResult.ExitVelocity, paResult.ExitVelocity...)
	gameRes.PAResult.LaunchAngle = append(gameRes.PAResult.LaunchAngle, paResult.LaunchAngle...)
	gameRes.PAResult.SprayAngle = append(gameRes.PAResult.SprayAngle, paResult.SprayAngle...)
	gameRes.PAResult.EventType = append(gameRes.PAResult.EventType, paResult.EventType...)
	gameRes.PAResult.AwayScore = append(gameRes.PAResult.AwayScore, paResult.AwayScore...)
	gameRes.PAResult.HomeScore = append(gameRes.PAResult.HomeScore, paResult.HomeScore...)
	gameRes.PAResult.AtBatNumber = append(gameRes.PAResult.AtBatNumber, paResult.AtBatNumber...)
	gameRes.PAResult.Inning = append(gameRes.PAResult.Inning, paResult.Inning...)
	gameRes.PAResult.InningTopBot = append(gameRes.PAResult.InningTopBot, paResult.InningTopBot...)
	gameRes.PAResult.Outs = append(gameRes.PAResult.Outs, paResult.Outs...)
	gameRes.PAResult.On1b = append(gameRes.PAResult.On1b, paResult.On1b...)
	gameRes.PAResult.On2b = append(gameRes.PAResult.On2b, paResult.On2b...)
	gameRes.PAResult.On3b = append(gameRes.PAResult.On3b, paResult.On3b...)
	gameRes.PAResult.RunnerOn1bID = append(gameRes.PAResult.RunnerOn1bID, paResult.RunnerOn1bID...)
	gameRes.PAResult.RunnerOn2bID = append(gameRes.PAResult.RunnerOn2bID, paResult.RunnerOn2bID...)
	gameRes.PAResult.RunnerOn3bID = append(gameRes.PAResult.RunnerOn3bID, paResult.RunnerOn3bID...)
	for i := 0; i < len(paResult.PitcherId); i++ {
		gameRes.PAResult.RbiBatterID = append(gameRes.PAResult.RbiBatterID, rbiBatterID)
		gameRes.PAResult.RbiCount = append(gameRes.PAResult.RbiCount, rbiCount)
	}
}

func postGameResults(gameRes []models.GameResult, season int, db *sql.DB) {
	gameYear := season

	for _, result := range gameRes {
		err := poster.InsertGameResult(db, result.GameId, result.GamePk, result.JobId, gameYear, result)
		if err != nil {
			fmt.Println("Error inserting game result:", err)
		}
	}
}

func advance(oldBases []int, basesMoved int, batterID int, event string) ([]int, []int) {
	newBases := make([]int, 3)
	runs := []int{}

	// 1) Advance existing runners (right-to-left to avoid overwrite)
	for i := 2; i >= 0; i-- {
		if oldBases[i] == 0 {
			continue
		}

		// sequencing boosts (optional house rules)
		if event == "single" && i == 1 && basesMoved == 1 { // runner on 2B scores on single
			runs = append(runs, oldBases[i])
			continue
		}
		if event == "double" && i == 0 && basesMoved == 2 { // runner on 1B scores on double
			runs = append(runs, oldBases[i])
			continue
		}

		dest := i + basesMoved
		if dest >= 3 {
			runs = append(runs, oldBases[i])
		} else {
			newBases[dest] = oldBases[i]
		}
	}

	// 2) Place the batter
	if event == "home_run" {
		runs = append(runs, batterID)
	} else {
		placement := basesMoved - 1 // 0 for walk/single, 1 for double, 2 for triple
		if placement < 0 {
			placement = 0
		}

		// force-advance if base occupied (walk with traffic, etc.)
		for placement < 3 && newBases[placement] != 0 {
			placement++
		}
		if placement >= 3 {
			runs = append(runs, batterID)
		} else {
			newBases[placement] = batterID
		}
	}

	return runs, newBases
}

// walkAdvance moves only the runners that are actually forced.
// returns (runs_scored, new_base_state)
func walkAdvance(old []int, batterID int) ([]int, []int) {
	newB := make([]int, 3)
	runs := []int{}

	// 3️⃣ Runner on 3B
	if old[0] != 0 && old[1] != 0 && old[2] != 0 {
		runs = append(runs, old[2])
	} else if old[2] != 0 {
		newB[2] = old[2]
	}

	// 2️⃣ Runner on 2B
	if old[1] != 0 {
		if old[0] != 0 { // 1B taken → forced
			newB[2] = old[1]
		} else {
			newB[1] = old[1]
		}
	}

	// 1️⃣ Runner on 1B
	if old[0] != 0 {
		newB[1] = old[0]
	}

	// Batter to 1B
	newB[0] = batterID

	return runs, newB
}

func ProcessPlateAppearance(
	paResult []models.PlateAppearanceResult,
	score int,
	baseState []int, // holds runner IDs
	outs int,
	gamePitchCount int,
) (int, []int, int, int, []int, int) {

	if len(paResult) == 0 || len(paResult[0].EventType) == 0 {
		return score, baseState, outs, gamePitchCount, nil, 0
	}
	lastEvent := paResult[0].EventType[len(paResult[0].EventType)-1]
	batterID := paResult[0].BatterId[len(paResult[0].BatterId)-1]
	scored := []int{}

	switch lastEvent {

	case "walk":
		var r []int
		r, baseState = walkAdvance(baseState, batterID)
		scored = append(scored, r...)

	case "single":
		var r []int
		r, baseState = advance(baseState, 1, batterID, "single")
		scored = append(scored, r...)

	case "double":
		var r []int
		r, baseState = advance(baseState, 2, batterID, "double")
		scored = append(scored, r...)

	case "triple":
		var r []int
		r, baseState = advance(baseState, 3, batterID, "triple")
		scored = append(scored, r...)
		// batter already placed on 3B by advance()

	case "home_run":
		for _, id := range baseState {
			if id != 0 {
				scored = append(scored, id)
			}
		}
		scored = append(scored, batterID)
		baseState = []int{0, 0, 0}

	case "out", "strikeout":
		outs++
	}
	score += len(scored)
	rbi := len(scored)

	updatedGamePitchCount := gamePitchCount + paResult[0].PitchCount[len(paResult[0].PitchCount)-1]

	return score, baseState, outs, updatedGamePitchCount, scored, rbi
}
