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

func SimulateGame(in []models.GameData) {

	db := config.ConnectDB()
	defer db.Close()

	homeLineup := []int{
		in[0].HomeBatter1Id, in[0].HomeBatter2Id, in[0].HomeBatter3Id,
		in[0].HomeBatter4Id, in[0].HomeBatter5Id, in[0].HomeBatter6Id,
		in[0].HomeBatter7Id, in[0].HomeBatter8Id, in[0].HomeBatter9Id,
	}

	homeLineupGameYear := []int{
		in[0].HomeBatter1GameYear, in[0].HomeBatter2GameYear, in[0].HomeBatter3GameYear,
		in[0].HomeBatter4GameYear, in[0].HomeBatter5GameYear, in[0].HomeBatter6GameYear,
		in[0].HomeBatter7GameYear, in[0].HomeBatter8GameYear, in[0].HomeBatter9GameYear,
	}

	awayLineup := []int{
		in[0].AwayBatter1Id, in[0].AwayBatter2Id, in[0].AwayBatter3Id,
		in[0].AwayBatter4Id, in[0].AwayBatter5Id, in[0].AwayBatter6Id,
		in[0].AwayBatter7Id, in[0].AwayBatter8Id, in[0].AwayBatter9Id,
	}

	awayLineupGameYear := []int{
		in[0].AwayBatter1GameYear, in[0].AwayBatter2GameYear, in[0].AwayBatter3GameYear,
		in[0].AwayBatter4GameYear, in[0].AwayBatter5GameYear, in[0].AwayBatter6GameYear,
		in[0].AwayBatter7GameYear, in[0].AwayBatter8GameYear, in[0].AwayBatter9GameYear,
	}

	homePitcherLineup := [][]int{
		{in[0].HomeStartingPitcherId, in[0].HomeStartingPitcherGameYear},
		{in[0].HomeBullpen1Id, in[0].HomeBullpen1GameYear},
		{in[0].HomeBullpen2Id, in[0].HomeBullpen2GameYear},
		{in[0].HomeBullpen3Id, in[0].HomeBullpen3GameYear},
		{in[0].HomeBullpen4Id, in[0].HomeBullpen4GameYear},
		{in[0].HomeBullpen5Id, in[0].HomeBullpen5GameYear},
		{in[0].HomeBullpen6Id, in[0].HomeBullpen6GameYear},
		{in[0].HomeBullpen7Id, in[0].HomeBullpen7GameYear},
		{in[0].HomeBullpen8Id, in[0].HomeBullpen8GameYear},
		{in[0].HomeBullpen9Id, in[0].HomeBullpen9GameYear},
	}

	awayPitcherLineup := [][]int{
		{in[0].AwayStartingPitcherId, in[0].AwayStartingPitcherGameYear},
		{in[0].AwayBullpen1Id, in[0].AwayBullpen1GameYear},
		{in[0].AwayBullpen2Id, in[0].AwayBullpen2GameYear},
		{in[0].AwayBullpen3Id, in[0].AwayBullpen3GameYear},
		{in[0].AwayBullpen4Id, in[0].AwayBullpen4GameYear},
		{in[0].AwayBullpen5Id, in[0].AwayBullpen5GameYear},
		{in[0].AwayBullpen6Id, in[0].AwayBullpen6GameYear},
		{in[0].AwayBullpen7Id, in[0].AwayBullpen7GameYear},
		{in[0].AwayBullpen8Id, in[0].AwayBullpen8GameYear},
		{in[0].AwayBullpen9Id, in[0].AwayBullpen9GameYear},
	}

	var homePitcher int
	var homePitcherGameYear int
	var awayPitcher int
	var awayPitcherGameYear int

	inning := 9
	awayScore := 0
	homeScore := 10
	awayBatterNumber := 0
	homeBatterNumber := 0
	atBatNumber := 0
	homePitcher = homePitcherLineup[0][0]
	homePitcherGameYear = homePitcherLineup[0][1]
	awayPitcher = awayPitcherLineup[0][0]
	awayPitcherGameYear = awayPitcherLineup[0][1]

	pitchingSubProbs, _ := fetcher.FetchPitchingSubstitutionProbs(db)

	for {

		topOuts := 0
		botOuts := 0
		awayBaseState := []bool{false, false, false, false}
		homeBaseState := []bool{false, false, false, false}
		priorAwayScore := awayScore
		priorHomeScore := homeScore
		gameRes := models.GameResult{
			GameId: uuid.New().String(),
			PAResult: models.PlateAppearanceResult{
				PitcherGameYear: []int{},
				PitcherFullName: []string{},
				PitcherId:       []int{},
				BatterGameYear:  []int{},
				BatterFullName:  []string{},
				BatterId:        []int{},
				BatterStands:    []string{},
				PitcherThrows:   []string{},
				Strikes:         []int{},
				Balls:           []int{},
				PitchCount:      []int{},
				PitchType:       []string{},
				PlateX:          []float64{},
				PlateZ:          []float64{},
				Zone:            []int{},
				Velocity:        []float64{},
				IsStrike:        []bool{},
				IsSwing:         []bool{},
				IsContact:       []string{},
				ExitVelocity:    []float64{},
				LaunchAngle:     []float64{},
				SprayAngle:      []float64{},
				EventType:       []string{},
				AwayScore:       []int{},
				HomeScore:       []int{},
				AtBatNumber:     []int{},
				Inning:          []int{},
				InningTopBot:    []string{},
				Outs:            []int{},
				On1b:            []bool{},
				On2b:            []bool{},
				On3b:            []bool{},
			}}

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
			}})

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
				homePitcherLineup = utils.FilterSliceSlices(homePitcherLineup, homePitcher)
				homePitcherChosenIndex := rand.Intn(len(homePitcherLineup))
				homePitcher = homePitcherLineup[homePitcherChosenIndex][0]
				homePitcherGameYear = homePitcherLineup[homePitcherChosenIndex][1]
			}

			fmt.Println("Batter #:", awayBatterNumber,
				"| Pitcher :", homePitcher,
				// "| Event:", awayPaResult[0].EventType[0],
				"| Event:", awayPaResult[0].EventType[len(awayPaResult[0].EventType)-1],
				"| EV:", awayPaResult[0].ExitVelocity[0],
				"| LA:", awayPaResult[0].LaunchAngle[0],
				"| SA:", awayPaResult[0].SprayAngle[0],
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
			}})

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
				awayPitcherLineup = utils.FilterSliceSlices(awayPitcherLineup, awayPitcher)
				awayPitcherChosenIndex := rand.Intn(len(awayPitcherLineup))
				homePitcher = awayPitcherLineup[awayPitcherChosenIndex][0]
				awayPitcherGameYear = awayPitcherLineup[awayPitcherChosenIndex][1]
			}

			fmt.Println("Batter #:", homeBatterNumber,
				"| Pitcher :", awayPitcher,
				// "| Event:", homePaResult[0].EventType[0],
				"| Event:", homePaResult[0].EventType[len(homePaResult[0].EventType)-1],
				"| EV:", homePaResult[0].ExitVelocity[0],
				"| LA:", homePaResult[0].LaunchAngle[0],
				"| SA:", homePaResult[0].SprayAngle[0],
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
