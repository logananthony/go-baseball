package sim

import (
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/logananthony/go-baseball/pkg/config"
	"github.com/logananthony/go-baseball/pkg/fetcher"
	"github.com/logananthony/go-baseball/pkg/models"
	"github.com/logananthony/go-baseball/pkg/poster"
	"github.com/logananthony/go-baseball/pkg/utils"
)

func SimulateGame(gameData []models.GameData) {

	db := config.ConnectDB()
	defer db.Close()

	simData := models.SimData{
		PlayerInfo:          []models.MLBPlayerInfo{},
		LeagueSwing:         []models.BatterSwingPercentageLeague{},
		LeagueContact:       []models.BatterContactPercentageLeague{},
		LeaguePitchCovMeans: []models.PitcherCovarianceMeanLeague{},
		BatterSwing:         []models.BatterSwingPercentage{},
		BatterContact:       []models.BatterContactPercentage{},
		BatterHitType:       []models.BatterHitType{},
		PitcherPitchFreq:    []models.PitcherCountPitchFreq{},
		PitcherCovMeans:     []models.PitcherCovarianceMean{},
		BatterEVDist:        []models.EVDistribution{},
		BatterLADist:        []models.LADistribution{},
		BatterSprayDist:     []models.SprayDistribution{},
	}

	gameRes := models.GameResult{
		GameId: uuid.New().String(),
	}

	homeBullpen := fetcher.FetchBullpenOrder(gameData[0].HomeTeam)
	awayBullpen := fetcher.FetchBullpenOrder(gameData[0].AwayTeam)

	// === FIX: Fetch starter pitcher info early and append to simData.PlayerInfo ===
	awayStarterInfoSlice, _ := fetcher.FetchPlayerInfo(db, &gameData[0].AwayStartingPitcher, &gameData[0].GameYear)
	homeStarterInfoSlice, _ := fetcher.FetchPlayerInfo(db, &gameData[0].HomeStartingPitcher, &gameData[0].GameYear)
	simData.PlayerInfo = append(simData.PlayerInfo, awayStarterInfoSlice...)
	simData.PlayerInfo = append(simData.PlayerInfo, homeStarterInfoSlice...)

	// Now safe to build playerInfoMap
	playerInfoMap := make(map[int]models.MLBPlayerInfo)
	for _, player := range simData.PlayerInfo {
		playerInfoMap[*player.ID] = player
	}

	// Lookup starter info safely
	awayStarterInfo := playerInfoMap[gameData[0].AwayStartingPitcher]
	awayStarterThrows := utils.ConvertPitcherThrows(awayStarterInfo.PitchHand)

	homeStarterInfo := playerInfoMap[gameData[0].HomeStartingPitcher]
	homeStarterThrows := utils.ConvertPitcherThrows(homeStarterInfo.PitchHand)

	// Fetch batting orders based on starter handedness
	homeBattingOrder, err := fetcher.FetchBattingOrder(gameData[0].HomeTeam, *awayStarterThrows)
	if err != nil {
		panic(err) // or return if wrapped in an API
	}

	awayBattingOrder, err := fetcher.FetchBattingOrder(gameData[0].AwayTeam, *homeStarterThrows)
	if err != nil {
		panic(err) // or return if wrapped in an API
	}

	// awayBattingOrder, _ := fetcher.FetchBattingOrder(gameData[0].AwayTeam, *homeStarterThrows)
	// homeBattingOrder, _ := fetcher.FetchBattingOrder(gameData[0].HomeTeam, *awayStarterThrows)

	homeLineup := []int{
		homeBattingOrder.PlayerID1,
		homeBattingOrder.PlayerID2,
		homeBattingOrder.PlayerID3,
		homeBattingOrder.PlayerID4,
		homeBattingOrder.PlayerID5,
		homeBattingOrder.PlayerID6,
		homeBattingOrder.PlayerID7,
		homeBattingOrder.PlayerID8,
		homeBattingOrder.PlayerID9,
	}

	for i, id := range homeLineup {
		if id == 0 {
			panic(fmt.Sprintf("Missing player in home batting order at slot %d — player ID is 0", i+1))
		}
	}

	homeLineupGameYear := make([]int, 9)
	for i := range homeLineupGameYear {
		homeLineupGameYear[i] = gameData[0].GameYear
	}

	awayLineup := []int{
		awayBattingOrder.PlayerID1,
		awayBattingOrder.PlayerID2,
		awayBattingOrder.PlayerID3,
		awayBattingOrder.PlayerID4,
		awayBattingOrder.PlayerID5,
		awayBattingOrder.PlayerID6,
		awayBattingOrder.PlayerID7,
		awayBattingOrder.PlayerID8,
		awayBattingOrder.PlayerID9,
	}

	for i, id := range awayLineup {
		if id == 0 {
			panic(fmt.Sprintf("Missing player in away batting order at slot %d — player ID is 0", i+1))
		}
	}

	awayLineupGameYear := make([]int, 9)
	for i := range awayLineupGameYear {
		awayLineupGameYear[i] = gameData[0].GameYear
	}

	homePitcherLineup := [][]int{
		{gameData[0].HomeStartingPitcher, gameData[0].GameYear},
		{homeBullpen.PlayerID1, gameData[0].GameYear},
		{homeBullpen.PlayerID2, gameData[0].GameYear},
		{homeBullpen.PlayerID3, gameData[0].GameYear},
		{homeBullpen.PlayerID4, gameData[0].GameYear},
		{homeBullpen.PlayerID5, gameData[0].GameYear},
		{homeBullpen.PlayerID6, gameData[0].GameYear},
		{homeBullpen.PlayerID7, gameData[0].GameYear},
		{homeBullpen.PlayerID8, gameData[0].GameYear},
	}

	awayPitcherLineup := [][]int{
		{gameData[0].AwayStartingPitcher, gameData[0].GameYear},
		{awayBullpen.PlayerID1, gameData[0].GameYear},
		{awayBullpen.PlayerID2, gameData[0].GameYear},
		{awayBullpen.PlayerID3, gameData[0].GameYear},
		{awayBullpen.PlayerID4, gameData[0].GameYear},
		{awayBullpen.PlayerID5, gameData[0].GameYear},
		{awayBullpen.PlayerID6, gameData[0].GameYear},
		{awayBullpen.PlayerID7, gameData[0].GameYear},
		{awayBullpen.PlayerID8, gameData[0].GameYear},
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

	fmt.Println("Caching data...")

	pitchingSubProbs, _ := fetcher.FetchPitchingSubstitutionProbs(db)
	simData.LeagueSwing = append(simData.LeagueSwing, fetcher.FetchBatterSwingPercentageLeague()...)
	simData.LeagueContact = append(simData.LeagueContact, fetcher.FetchBatterContactPercentageLeague()...)
	simData.LeaguePitchCovMeans = append(simData.LeaguePitchCovMeans, fetcher.FetchPitcherCovarianceMeanLeague()...)

	for i := 0; i <= 8; i++ {

		homeBatter := homeLineup[i]
		homeBatterGameYear := homeLineupGameYear[i]
		homeBatterInfo, _ := fetcher.FetchPlayerInfo(db, &homeBatter, &homeBatterGameYear)
		homeBatterSwingProbs, _ := fetcher.FetchBatterSwingPercentage(db, homeBatter, homeBatterGameYear)
		homeBatterContactProbs, _ := fetcher.FetchBatterContactPercentage(db, homeBatter, homeBatterGameYear)
		homeBatterHitProbs, _ := fetcher.FetchBatterHitType(db, homeBatter, homeBatterGameYear)
		homeBatterEvDist := fetcher.FetchEVDistributions(db, homeBatterGameYear, homeBatter)
		homeBatterLaDist := fetcher.FetchLADistributions(db, homeBatterGameYear, homeBatter)
		homeBatterSprayDist := fetcher.FetchSprayDistributions(db, homeBatterGameYear, homeBatter)

		simData.PlayerInfo = append(simData.PlayerInfo, homeBatterInfo...)
		simData.BatterSwing = append(simData.BatterSwing, homeBatterSwingProbs...)
		simData.BatterContact = append(simData.BatterContact, homeBatterContactProbs...)
		simData.BatterHitType = append(simData.BatterHitType, homeBatterHitProbs...)
		simData.BatterEVDist = append(simData.BatterEVDist, homeBatterEvDist...)
		simData.BatterLADist = append(simData.BatterLADist, homeBatterLaDist...)
		simData.BatterSprayDist = append(simData.BatterSprayDist, homeBatterSprayDist...)

		awayBatter := awayLineup[i]
		awayBatterGameYear := awayLineupGameYear[i]
		awayBatterInfo, _ := fetcher.FetchPlayerInfo(db, &awayBatter, &awayBatterGameYear)
		awayBatterSwingProbs, _ := fetcher.FetchBatterSwingPercentage(db, awayBatter, awayBatterGameYear)
		awayBatterContactProbs, _ := fetcher.FetchBatterContactPercentage(db, awayBatter, awayBatterGameYear)
		awayBatterHitProbs, _ := fetcher.FetchBatterHitType(db, awayBatter, awayBatterGameYear)
		awayBatterEvDist := fetcher.FetchEVDistributions(db, awayBatterGameYear, awayBatter)
		awayBatterLaDist := fetcher.FetchLADistributions(db, awayBatterGameYear, awayBatter)
		awayBatterSprayDist := fetcher.FetchSprayDistributions(db, awayBatterGameYear, awayBatter)

		simData.PlayerInfo = append(simData.PlayerInfo, awayBatterInfo...)
		simData.BatterSwing = append(simData.BatterSwing, awayBatterSwingProbs...)
		simData.BatterContact = append(simData.BatterContact, awayBatterContactProbs...)
		simData.BatterHitType = append(simData.BatterHitType, awayBatterHitProbs...)
		simData.BatterEVDist = append(simData.BatterEVDist, awayBatterEvDist...)
		simData.BatterLADist = append(simData.BatterLADist, awayBatterLaDist...)
		simData.BatterSprayDist = append(simData.BatterSprayDist, awayBatterSprayDist...)
	}

	for i := 0; i <= 8; i++ {

		homePitcher := homePitcherLineup[i][0]
		homePitcherGameYear := homePitcherLineup[i][1]
		homePitchFreqsR := fetcher.FetchPitcherFrequencies(db, homePitcher, "R")
		homePitchFreqsL := fetcher.FetchPitcherFrequencies(db, homePitcher, "L")
		homePitcherCov := fetcher.FetchPitcherCovarianceMean(db, int64(homePitcher), int64(homePitcherGameYear))
		homePitcherInfo, _ := fetcher.FetchPlayerInfo(db, &homePitcher, &homePitcherGameYear)

		simData.PitcherPitchFreq = append(simData.PitcherPitchFreq, homePitchFreqsR...)
		simData.PitcherPitchFreq = append(simData.PitcherPitchFreq, homePitchFreqsL...)
		simData.PitcherCovMeans = append(simData.PitcherCovMeans, homePitcherCov...)
		simData.PlayerInfo = append(simData.PlayerInfo, homePitcherInfo...)

		awayPitcher := awayPitcherLineup[i][0]
		awayPitcherGameYear := awayPitcherLineup[i][1]
		awayPitchFreqsR := fetcher.FetchPitcherFrequencies(db, awayPitcher, "R")
		awayPitchFreqsL := fetcher.FetchPitcherFrequencies(db, awayPitcher, "L")
		awayPitcherCov := fetcher.FetchPitcherCovarianceMean(db, int64(awayPitcher), int64(awayPitcherGameYear))
		awayPitcherInfo, _ := fetcher.FetchPlayerInfo(db, &awayPitcher, &awayPitcherGameYear)

		simData.PitcherPitchFreq = append(simData.PitcherPitchFreq, awayPitchFreqsR...)
		simData.PitcherPitchFreq = append(simData.PitcherPitchFreq, awayPitchFreqsL...)
		simData.PitcherCovMeans = append(simData.PitcherCovMeans, awayPitcherCov...)
		simData.PlayerInfo = append(simData.PlayerInfo, awayPitcherInfo...)
	}

	fmt.Println("simData has been written to sim_data.json")

	fmt.Println("Done caching data.")

	for {

		topOuts := 0
		botOuts := 0
		awayBaseState := []bool{false, false, false, false}
		homeBaseState := []bool{false, false, false, false}
		priorAwayScore := awayScore
		priorHomeScore := homeScore

		fmt.Println("Top Inning:", inning)

		for topOuts < 3 {
			awayBatterNumber = awayBatterNumber % 9
			awayBatter := awayLineup[awayBatterNumber]
			awayBatterGameYear := awayLineupGameYear[awayBatterNumber]
			awayPaResult := SimulatePlateAppearance([]models.PlateAppearanceData{{
				BatterGameYear:  awayBatterGameYear,
				BatterId:        awayBatter,
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

			// spew.Dump(awayPaResult)

			atBatNumber++

			inningRunsHome := awayScore - priorAwayScore
			pullProbHome := utils.GetPullProbability(pitchingSubProbs, inning, awayScore, inningRunsHome)

			var pitcherPulledHome bool

			if pullProbHome != nil {
				pitcherPulledHome = utils.IsSuccess(pullProbHome)
			} else {
				pitcherPulledHome = false
			}

			if pullProbHome != nil && pitcherPulledHome {
				if len(homePitcherLineup) > 0 {
					homePitcherLineup = utils.FilterSliceSlices(homePitcherLineup, homePitcher)
					if len(homePitcherLineup) > 0 { // Ensure the lineup is not empty after filtering
						homePitcherChosenIndex := rand.Intn(len(homePitcherLineup))
						homePitcher = homePitcherLineup[homePitcherChosenIndex][0]
						homePitcherGameYear = homePitcherLineup[homePitcherChosenIndex][1]
					} else {
						fmt.Println("Home pitcher lineup is empty after filtering, skipping pitcher substitution.")
					}
				} else {
					fmt.Println("Home pitcher lineup is empty, skipping pitcher substitution.")
				}
			}

			fmt.Println("Batter #:", awayBatterNumber,
				"| Pitcher :", homePitcher,
				"| Event:", awayPaResult[0].EventType[len(awayPaResult[0].EventType)-1],
				"| EV:", awayPaResult[0].ExitVelocity[len(awayPaResult[0].ExitVelocity)-1],
				"| LA:", awayPaResult[0].LaunchAngle[len(awayPaResult[0].LaunchAngle)-1],
				"| SA:", awayPaResult[0].SprayAngle[len(awayPaResult[0].SprayAngle)-1],
				"| Base State:", awayBaseState[0], awayBaseState[1], awayBaseState[2],
				"| Score:", awayScore, "-", homeScore)

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
			postGameResults([]models.GameResult{gameRes}, db)
			fmt.Println("Home team wins:", homeScore, "-", awayScore)
			break
		}

		fmt.Println("Bottom Inning:", inning)

		for botOuts < 3 {
			homeBatterNumber = homeBatterNumber % 9
			homeBatter := homeLineup[homeBatterNumber]
			homeBatterGameYear := homeLineupGameYear[homeBatterNumber]
			homePaResult := SimulatePlateAppearance([]models.PlateAppearanceData{{
				BatterGameYear:  homeBatterGameYear,
				BatterId:        homeBatter,
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

			inningRunsAway := homeScore - priorHomeScore
			pullProbAway := utils.GetPullProbability(pitchingSubProbs, inning, homeScore, inningRunsAway)

			var pitcherPulledAway bool

			if pullProbAway != nil {
				pitcherPulledAway = utils.IsSuccess(pullProbAway)
			} else {
				pitcherPulledAway = false
			}

			if pullProbAway != nil && pitcherPulledAway {
				if len(awayPitcherLineup) > 0 {
					awayPitcherLineup = utils.FilterSliceSlices(awayPitcherLineup, awayPitcher)
					if len(awayPitcherLineup) > 0 { // Ensure the lineup is not empty after filtering
						awayPitcherChosenIndex := rand.Intn(len(awayPitcherLineup))
						awayPitcher = awayPitcherLineup[awayPitcherChosenIndex][0]
						awayPitcherGameYear = awayPitcherLineup[awayPitcherChosenIndex][1]
					} else {
						fmt.Println("Away pitcher lineup is empty after filtering, skipping pitcher substitution.")
					}
				} else {
					fmt.Println("Away pitcher lineup is empty, skipping pitcher substitution.")
				}
			}

			fmt.Println("Batter #:", homeBatterNumber,
				"| Pitcher :", awayPitcher,
				"| Event:", homePaResult[0].EventType[len(homePaResult[0].EventType)-1],
				"| EV:", homePaResult[0].ExitVelocity[len(homePaResult[0].ExitVelocity)-1],
				"| LA:", homePaResult[0].LaunchAngle[len(homePaResult[0].LaunchAngle)-1],
				"| SA:", homePaResult[0].SprayAngle[len(homePaResult[0].SprayAngle)-1],
				"| Base State:", homeBaseState[0], homeBaseState[1], homeBaseState[2],
				"| Score:", awayScore, "-", homeScore)

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
				postGameResults([]models.GameResult{gameRes}, db)
			}
		}

		// If 9 or later and not tied, game ends
		if inning >= 9 && homeScore != awayScore {
			fmt.Println("Away team wins:", awayScore, "-", homeScore)
			postGameResults([]models.GameResult{gameRes}, db)
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

// func AppendSimData(simData *models.SimData,
// 	playerInfo models.MLBPlayerInfo,
// 	leagueSwing models.BatterSwingPercentageLeague,
// 	leagueContact models.BatterContactPercentageLeague,
// 	leaguePitchCovMeans models.PitcherCovarianceMeanLeague,
// 	batterSwing models.BatterSwingPercentage,
// 	batterContact models.BatterContactPercentage,
// 	batterHitType models.BatterHitType,
// 	pitcherPitchFreq models.PitcherCountPitchFreq,
// 	pitcherCovMeans models.PitcherCovarianceMean,
// 	batterEVDist models.EVDistribution,
// 	batterLADist models.LADistribution,
// 	batterSprayDist models.SprayDistribution) {
// 	simData.PlayerInfo = append(simData.PlayerInfo, playerInfo)
// 	simData.LeagueSwing = append(simData.LeagueSwing, leagueSwing)
// 	simData.LeagueContact = append(simData.LeagueContact, leagueContact)
// 	simData.LeaguePitchCovMeans = append(simData.LeaguePitchCovMeans, leaguePitchCovMeans)
// 	simData.BatterSwing = append(simData.BatterSwing, batterSwing)
// 	simData.BatterContact = append(simData.BatterContact, batterContact)
// 	simData.BatterHitType = append(simData.BatterHitType, batterHitType)
// 	simData.PitcherPitchFreq = append(simData.PitcherPitchFreq, pitcherPitchFreq)
// 	simData.PitcherCovMeans = append(simData.PitcherCovMeans, pitcherCovMeans)
// 	simData.BatterEVDist = append(simData.BatterEVDist, batterEVDist)
// 	simData.BatterLADist = append(simData.BatterLADist, batterLADist)
// 	simData.BatterSprayDist = append(simData.BatterSprayDist, batterSprayDist)
// }

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

func ProcessPlateAppearance(paResult []models.PlateAppearanceResult, score int, baseState []bool, outs int) (int, []bool, int) {

	db := config.ConnectDB()
	defer db.Close()

	if len(paResult) == 0 || len(paResult[0].EventType) == 0 {
		return score, baseState, outs
	}

	lastEventIndex := len(paResult[0].EventType) - 1
	lastEventType := paResult[0].EventType[lastEventIndex]

	switch lastEventType {
	case "walk":
		if baseState[0] && baseState[1] && baseState[2] {
			baseState[3] = true
			baseState[2] = true
			baseState[1] = true
			baseState[0] = true
		} else {
			if baseState[1] && baseState[2] {
				baseState[2] = true
				baseState[1] = true
				baseState[0] = true
			} else if baseState[0] && baseState[2] {
				baseState[1] = baseState[0]
				baseState[0] = true
			} else if baseState[0] && baseState[1] {
				baseState[1] = baseState[0]
				baseState[2] = baseState[1]
				baseState[0] = true
			} else if baseState[2] {
				baseState[0] = true
			} else if baseState[1] {
				baseState[0] = true
			} else if baseState[0] {
				baseState[1] = baseState[0]
				baseState[0] = true
			} else {
				baseState[0] = true
			}
		}

	case "single", "double", "triple":

		priorBaseState := baseState
		newBaseState := [4]bool{false, false, false, false}

		for i := range baseState {

			if priorBaseState[i] {
				var basesMoved int
				if lastEventType == "single" {
					basesMoved = 1
				} else if lastEventType == "double" {
					basesMoved = 2
				} else {
					basesMoved = 3
				}

				if i+basesMoved < 3 {
					newBaseState[i+basesMoved] = true
				} else {
					score++
				}

			} else {
				if lastEventType == "single" {
					newBaseState[0] = true
				} else if lastEventType == "double" {
					newBaseState[1] = true
				} else {
					newBaseState[2] = true
				}
				baseState = newBaseState[:]
				break
			}
			baseState = newBaseState[:]
		}

	case "home_run":
		runs := 1
		for i := 0; i <= 2; i++ {
			if baseState[i] {
				runs++
				baseState[i] = false
			}
		}
		score += runs
	case "out", "strikeout":
		outs++
	}

	if baseState[3] {
		score++
		baseState[3] = false
	}

	return score, baseState, outs
}

func postGameResults(gameRes []models.GameResult, db *sql.DB) {
	gameYear := 2024

	for _, result := range gameRes {
		err := poster.InsertGameResult(db, result.GameId, gameYear, result)
		if err != nil {
			fmt.Println("Error inserting game result:", err)
		}
	}
}
