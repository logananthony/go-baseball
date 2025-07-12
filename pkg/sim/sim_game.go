package sim

import (
	"database/sql"
	"fmt"
	"log"

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
	gameData models.GameData, // for ids
) {

	// fmt.Printf("simData.PlayerInfo count: %d\n", len(simData.PlayerInfo))
	// fmt.Printf("simData.BatterSwing count: %d\n", len(simData.BatterSwing))
	// fmt.Printf("simData.BatterContact count: %d\n", len(simData.BatterContact))
	// fmt.Printf("simData.BatterHitType count: %d\n", len(simData.BatterHitType))
	// fmt.Printf("simData.LeagueContact count: %d\n", len(simData.LeagueContact))
	// Add any others that are relevant for your lookups

	defer func() {
		if r := recover(); r != nil {
			log.Printf("🔥 Panic recovered in SimulateGame: %v", r)
		}
	}()

	// db := config.ConnectDB()
	// defer db.Close()

	gamePk, _ := fetcher.FetchGameDataByGamePk(db, gameData.GamePk)

	first := gamePk[0]

	homeStartingPitcher := first.HomePitcherId
	awayStartingPitcher := first.AwayPitcherId
	homeTeam := first.HomeTeamAbbr
	awayTeam := first.AwayTeamAbbr
	season := first.Season

	gameRes := models.GameResult{
		GameId: uuid.New().String(),
		GamePk: gameData.GamePk,
		JobId:  gameData.JobId,
	}

	// === FIX: Fetch starter pitcher info early and append to simData.PlayerInfo ===
	// awayStarterInfoSlice, _ := fetcher.FetchPlayerInfo(db, awayStartingPitcher)
	// homeStarterInfoSlice, _ := fetcher.FetchPlayerInfo(db, homeStartingPitcher)
	// simData.PlayerInfo = append(simData.PlayerInfo, awayStarterInfoSlice...)
	// simData.PlayerInfo = append(simData.PlayerInfo, homeStarterInfoSlice...)

	// fmt.Println("Fetching full game lineup data...")
	// allGameData, _ := fetcher.FetchGameDataByGamePk(db, gameData.GamePk)
	// fmt.Println("Fetched full game lineup data")

	// for _, entry := range allGameData {
	// 	switch entry.Team {
	// 	case "home":
	// 		homeLineup = append(homeLineup, entry)
	// 	case "away":
	// 		awayLineup = append(awayLineup, entry)
	// 	}
	// }

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
	awayScore := 0
	homeScore := 0
	awayBatterNumber := 0
	homeBatterNumber := 0
	atBatNumber := 0
	homePitcher = homePitcherLineup[0][0]
	homePitcherGameYear = homePitcherLineup[0][1]
	awayPitcher = awayPitcherLineup[0][0]
	awayPitcherGameYear = awayPitcherLineup[0][1]

	// fmt.Println("Caching league-level data...")

	simData.HomeTeam = homeTeam

	for {

		topOuts := 0
		botOuts := 0
		awayBaseState := []bool{false, false, false, false}
		homeBaseState := []bool{false, false, false, false}
		priorAwayScore := awayScore
		priorHomeScore := homeScore

		inningRunsHome := awayScore - priorAwayScore
		pullProbHome := utils.GetPullProbability(pitchingSubProbs, inning, awayScore, inningRunsHome)

		var pitcherPulledHome bool

		if pullProbHome != nil {
			pitcherPulledHome = utils.IsSuccess(pullProbHome)
		} else {
			pitcherPulledHome = false
		}
		// fmt.Printf("✅ Cached %d pitcher substitution probabilities\n", len(pitchingSubProbs))
		// for _, p := range pitchingSubProbs[:min(3, len(pitchingSubProbs))] {
		// 	fmt.Printf("Example substitution rule: %+v\n", p)
		// }

		usedPitchersHome := map[int]bool{}

		if pullProbHome != nil && pitcherPulledHome {
			if len(homePitcherLineup) > 0 {
				homePitcherLineup = utils.FilterSliceSlices(homePitcherLineup, homePitcher)
				if len(homePitcherLineup) > 0 {
					// homePitcherChosenIndex := rand.Intn(len(homePitcherLineup))
					// homePitcher = homePitcherLineup[homePitcherChosenIndex][0]
					// homePitcherGameYear = homePitcherLineup[homePitcherChosenIndex][1]
					runDiff := awayScore - homeScore // for homePitcher (home is pitching)
					runnersOn := utils.CountRunners(awayBaseState)

					selected := utils.SelectBullpenPitcherLineup(
						homePitcherLineup,
						bullpenRoleProbs,
						inning, runDiff, runnersOn, usedPitchersHome,
					)

					if selected != nil {
						homePitcher = selected[0]
						homePitcherGameYear = selected[1]
						usedPitchersHome[homePitcher] = true
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

		inningRunsAway := homeScore - priorHomeScore
		pullProbAway := utils.GetPullProbability(pitchingSubProbs, inning, homeScore, inningRunsAway)

		var pitcherPulledAway bool

		if pullProbAway != nil {
			pitcherPulledAway = utils.IsSuccess(pullProbAway)
		} else {
			pitcherPulledAway = false
		}

		usedPitchersAway := map[int]bool{}

		if pullProbAway != nil && pitcherPulledAway {
			if len(awayPitcherLineup) > 0 {
				awayPitcherLineup = utils.FilterSliceSlices(awayPitcherLineup, awayPitcher)
				if len(awayPitcherLineup) > 0 { // Ensure the lineup is not empty after filtering
					// awayPitcherChosenIndex := rand.Intn(len(awayPitcherLineup))
					// awayPitcher = awayPitcherLineup[awayPitcherChosenIndex][0]
					// awayPitcherGameYear = awayPitcherLineup[awayPitcherChosenIndex][1]
					runDiff := homeScore - awayScore // for awayPitcher (away is pitching)
					runnersOn := utils.CountRunners(homeBaseState)

					selected := utils.SelectBullpenPitcherLineup(
						awayPitcherLineup,
						bullpenRoleProbs,
						inning, runDiff, runnersOn, usedPitchersAway,
					)
					if selected != nil {
						awayPitcher = selected[0]
						awayPitcherGameYear = selected[1]
						usedPitchersAway[awayPitcher] = true
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

		// fmt.Println("Top Inning:", inning)

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
				On1b:            awayBaseState[0],
				On2b:            awayBaseState[1],
				On3b:            awayBaseState[2],
			}}, []models.SimData{simData})

			atBatNumber++

			// ────────────────────────────────────────────────────────────────────────

			// fmt.Println("Batter #:", awayBatterNumber,
			// 	"| Pitcher :", homePitcher,
			// 	"| Event:", awayPaResult[0].EventType[len(awayPaResult[0].EventType)-1],
			// 	// "| EV:", awayPaResult[0].ExitVelocity[len(awayPaResult[0].ExitVelocity)-1],
			// 	// "| LA:", awayPaResult[0].LaunchAngle[len(awayPaResult[0].LaunchAngle)-1],
			// 	// "| SA:", awayPaResult[0].SprayAngle[len(awayPaResult[0].SprayAngle)-1],
			// 	"| Base State:", awayBaseState[0], awayBaseState[1], awayBaseState[2],
			// 	"| Score:", awayScore, "-", homeScore)

			awayScore, awayBaseState, topOuts = ProcessPlateAppearance(
				awayPaResult, awayScore, awayBaseState, topOuts,
			)

			for _, paResult := range awayPaResult {
				AppendPlateAppearanceTopResult(paResult, awayScore, homeScore, atBatNumber, inning, topOuts, awayBaseState)
				AppendGameResult(&gameRes, paResult)
			}
			// spew.Dump(gameRes.PAResult)

			awayBatterNumber++
		}

		if inning >= 9 && homeScore > awayScore {
			postGameResults([]models.GameResult{gameRes}, season, db)
			fmt.Println("Home team wins:", homeScore, "-", awayScore)
			break
		}

		// fmt.Println("Bottom Inning:", inning)

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
				On1b:            homeBaseState[0],
				On2b:            homeBaseState[1],
				On3b:            homeBaseState[2],
			}}, []models.SimData{simData})

			atBatNumber++

			// fmt.Println("Batter #:", homeBatterNumber,
			// 	"| Pitcher :", awayPitcher,
			// 	"| Event:", homePaResult[0].EventType[len(homePaResult[0].EventType)-1],
			// 	// "| EV:", homePaResult[0].ExitVelocity[len(homePaResult[0].ExitVelocity)-1],
			// 	// "| LA:", homePaResult[0].LaunchAngle[len(homePaResult[0].LaunchAngle)-1],
			// 	// "| SA:", homePaResult[0].SprayAngle[len(homePaResult[0].SprayAngle)-1],
			// 	"| Base State:", homeBaseState[0], homeBaseState[1], homeBaseState[2],
			// 	"| Score:", awayScore, "-", homeScore)

			homeScore, homeBaseState, botOuts = ProcessPlateAppearance(
				homePaResult, homeScore, homeBaseState, botOuts,
			)

			for _, paResult := range homePaResult {
				AppendPlateAppearanceBotResult(paResult, awayScore, homeScore, atBatNumber, inning, topOuts, homeBaseState)
				AppendGameResult(&gameRes, paResult)
			}
			// spew.Dump(gameRes.PAResult)
			homeBatterNumber++

			if inning >= 9 && homeScore > awayScore {
				fmt.Println("Home team wins (walk-off):", homeScore, "-", awayScore)
				postGameResults([]models.GameResult{gameRes}, season, db)
			}
		}

		// If 9 or later and not tied, game ends
		if inning >= 9 && homeScore != awayScore {
			fmt.Println("Away team wins:", awayScore, "-", homeScore)
			postGameResults([]models.GameResult{gameRes}, season, db)
			break
		}

		inning++
	}

}

func AppendPlateAppearanceTopResult(paResult models.PlateAppearanceResult, awayScore int, homeScore int, atBatNumber int, inning int, topOuts int, awayBaseState []bool) {
	paResult.AwayScore = append(paResult.AwayScore, awayScore)
	paResult.HomeScore = append(paResult.HomeScore, homeScore)
	paResult.AtBatNumber = append(paResult.AtBatNumber, atBatNumber)
	paResult.Inning = append(paResult.Inning, inning)
	paResult.InningTopBot = append(paResult.InningTopBot, "Top")
	paResult.Outs = append(paResult.Outs, topOuts)
	paResult.On1b = append(paResult.On1b, awayBaseState[0])
	paResult.On2b = append(paResult.On2b, awayBaseState[1])
	paResult.On3b = append(paResult.On3b, awayBaseState[2])
}

func AppendPlateAppearanceBotResult(paResult models.PlateAppearanceResult, awayScore int, homeScore int, atBatNumber int, inning int, topOuts int, homeBaseState []bool) {
	paResult.AwayScore = append(paResult.AwayScore, awayScore)
	paResult.HomeScore = append(paResult.HomeScore, homeScore)
	paResult.AtBatNumber = append(paResult.AtBatNumber, atBatNumber)
	paResult.Inning = append(paResult.Inning, inning)
	paResult.InningTopBot = append(paResult.InningTopBot, "Bot")
	paResult.Outs = append(paResult.Outs, topOuts)
	paResult.On1b = append(paResult.On1b, homeBaseState[0])
	paResult.On2b = append(paResult.On2b, homeBaseState[1])
	paResult.On3b = append(paResult.On3b, homeBaseState[2])
}

func AppendGameResult(gameRes *models.GameResult, paResult models.PlateAppearanceResult) {
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

func advance(oldBases []bool, basesMoved int, event string) (int, []bool) {
	newBases := make([]bool, 3)
	runs := 0

	// 1) Advance existing runners (right-to-left to avoid overwrite)
	for i := 2; i >= 0; i-- {
		if !oldBases[i] {
			continue
		}

		// sequencing boosts (optional house rules)
		if event == "single" && i == 1 && basesMoved == 1 { // runner on 2B scores on single
			runs++
			continue
		}
		if event == "double" && i == 0 && basesMoved == 2 { // runner on 1B scores on double
			runs++
			continue
		}

		dest := i + basesMoved
		if dest >= 3 {
			runs++ // runner scores
		} else {
			newBases[dest] = true
		}
	}

	// 2) Place the batter
	if event == "home_run" {
		runs++ // batter scores; bases cleared in caller
	} else {
		placement := basesMoved - 1 // 0 for walk/single, 1 for double, 2 for triple
		if placement < 0 {
			placement = 0
		}

		// force-advance if base occupied (walk with traffic, etc.)
		for placement < 3 && newBases[placement] {
			placement++
		}
		if placement >= 3 {
			runs++ // forced across home (e.g., walk with bases loaded)
		} else {
			newBases[placement] = true
		}
	}

	return runs, newBases
}

// walkAdvance moves only the runners that are actually forced.
// returns (runs_scored, new_base_state)
func walkAdvance(old []bool) (int, []bool) {
	newB := make([]bool, 3)
	runs := 0

	// 3️⃣ Runner on 3B
	if old[0] && old[1] && old[2] {
		// bases were loaded → he’s forced home
		runs++
	} else if old[2] {
		newB[2] = true // stays put
	}

	// 2️⃣ Runner on 2B
	if old[1] {
		if old[0] { // 1B taken → forced
			newB[2] = true // to 3B
		} else {
			newB[1] = true // stays
		}
	}

	// 1️⃣ Runner on 1B
	if old[0] {
		newB[1] = true // always goes to 2B
	}

	// Batter to 1B
	newB[0] = true

	return runs, newB
}

func ProcessPlateAppearance(
	paResult []models.PlateAppearanceResult,
	score int,
	baseState []bool, // must be len==3 now
	outs int,
) (int, []bool, int) {

	if len(paResult) == 0 || len(paResult[0].EventType) == 0 {
		return score, baseState, outs
	}
	lastEvent := paResult[0].EventType[len(paResult[0].EventType)-1]

	var runs int

	switch lastEvent {

	case "walk":
		var r int
		r, baseState = walkAdvance(baseState)
		score += r

	case "single":
		runs, baseState = advance(baseState, 1, "single")
		score += runs

	case "double":
		runs, baseState = advance(baseState, 2, "double")
		score += runs

	case "triple":
		runs, baseState = advance(baseState, 3, "triple")
		score += runs
		// batter already placed on 3B by advance()

	case "home_run":
		for _, occ := range baseState {
			if occ {
				score++
			}
		}
		score++ // batter
		baseState = []bool{false, false, false}

	case "out", "strikeout":
		outs++
	}

	// -----------------------------------------------------------------
	// (Optional) extra-base adjustments – uncomment & tune percentages
	/*
		if lastEvent == "single" {
			// runner on 2B scores on single 70 %
			if baseState[1] && rand.Float64() < 0.70 {
				baseState[1] = false
				score++
			}
			// runner on 1B advances to 3B 35 %
			if baseState[0] && rand.Float64() < 0.35 {
				baseState[0] = false
				baseState[2] = true
			}
		}
	*/
	// -----------------------------------------------------------------

	return score, baseState, outs
}
