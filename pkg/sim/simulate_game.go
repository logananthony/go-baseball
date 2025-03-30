package sim

import (
	"fmt"
	"github.com/logananthony/go-baseball/pkg/config"
	"github.com/logananthony/go-baseball/pkg/models"
)

func SimulateGame(in []models.GameData) []models.GameResult {
	db := config.ConnectDB()
	defer db.Close()

	homeLineup := []int{
		in[0].HomeBatter1, in[0].HomeBatter2, in[0].HomeBatter3,
		in[0].HomeBatter4, in[0].HomeBatter5, in[0].HomeBatter6,
		in[0].HomeBatter7, in[0].HomeBatter8, in[0].HomeBatter9,
	}

	awayLineup := []int{
		in[0].AwayBatter1, in[0].AwayBatter2, in[0].AwayBatter3,
		in[0].AwayBatter4, in[0].AwayBatter5, in[0].AwayBatter6,
		in[0].AwayBatter7, in[0].AwayBatter8, in[0].AwayBatter9,
	}

	homePitcher := in[0].HomeStartingPitcher
	awayPitcher := in[0].AwayStartingPitcher

	inning := 9
	awayScore := 0
	homeScore := 0
	awayBatterNumber := 0
	homeBatterNumber := 0
	atBatNumber := 0

	var gameRes []models.GameResult

	for {
		topOuts := 0
		botOuts := 0
		awayBaseState := []bool{false, false, false, false}
		homeBaseState := []bool{false, false, false, false}

		fmt.Println("Top Inning:", inning)

		// --- Top Half ---
		for topOuts < 3 {
      awayBatterNumber = awayBatterNumber % 9
			awayBatter := awayLineup[awayBatterNumber]
			awayPaResult := SimulateAtBat([]models.PlateAppearanceData{{
				GameYear: 2024,
				PitcherId: homePitcher,
				BatterId:  awayBatter,
				Strikes:   0,
				Balls:     0,
			}})

			atBatNumber++
			awayScore, awayBaseState, topOuts = ProcessPlateAppearance(
				awayPaResult, awayScore, awayBaseState, topOuts,
			)

			fmt.Println("Batter #:", awayBatterNumber, "| Event:", awayPaResult[0].EventType[0], "| Base State:", awayBaseState[0], awayBaseState[1], awayBaseState[2], "| Score:", awayScore, "-", homeScore)

			gameRes = append(gameRes, BuildGameResult(awayPaResult, atBatNumber, inning, "Top", topOuts, awayBaseState, awayScore, homeScore))
			awayBatterNumber++
		}

		// Early end if home team leads after top 9
		if inning >= 9 && homeScore > awayScore {
			fmt.Println("Home team wins:", homeScore, "-", awayScore)
			break
		}

		fmt.Println("Bottom Inning:", inning)

		for botOuts < 3 {
      homeBatterNumber = homeBatterNumber % 9
			homeBatter := homeLineup[homeBatterNumber]
			homePaResult := SimulateAtBat([]models.PlateAppearanceData{{
				GameYear: 2024,
				PitcherId: awayPitcher,
				BatterId:  homeBatter,
				Strikes:   0,
				Balls:     0,
			}})

			atBatNumber++
			homeScore, homeBaseState, botOuts = ProcessPlateAppearance(
				homePaResult, homeScore, homeBaseState, botOuts,
			)

			fmt.Println("Batter #:", homeBatterNumber, "| Event:", homePaResult[0].EventType[0], "| Base State:", homeBaseState[0], homeBaseState[1], homeBaseState[2], "| Score:", awayScore, "-", homeScore)

			gameRes = append(gameRes, BuildGameResult(homePaResult, atBatNumber, inning, "Bot", botOuts, homeBaseState, awayScore, homeScore))
			homeBatterNumber++

			// Walk-off condition
			if inning >= 9 && homeScore > awayScore {
				fmt.Println("Home team wins (walk-off):", homeScore, "-", awayScore)
				return gameRes
			}
		}

		// If 9 or later and not tied, game ends
		if inning >= 9 && homeScore != awayScore {
			fmt.Println("Away team wins:", awayScore, "-", homeScore)
			break
		}

		inning++
	}

	return gameRes
}

func ProcessPlateAppearance(paResult []models.PlateAppearanceResult, score int, baseState []bool, outs int) (int, []bool, int) {
	if len(paResult) == 0 || len(paResult[0].EventType) == 0 {
		return score, baseState, outs
	}

	switch paResult[0].EventType[0] {
	case "walk":
		if !baseState[0] {
			baseState[0] = true
		}
	case "single":
		baseState[3] = baseState[2]
		baseState[2] = baseState[1]
		baseState[1] = baseState[0]
		baseState[0] = true
	case "double":
		baseState[3] = baseState[1]
		baseState[2] = baseState[0]
		baseState[1] = true
		baseState[0] = false
	case "triple":
		baseState[3] = baseState[0] || baseState[1] || baseState[2]
		baseState[2] = true
		baseState[1] = false
		baseState[0] = false
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

func BuildGameResult(pa []models.PlateAppearanceResult, abNum, inning int, half string, outs int, bases []bool, awayScore, homeScore int) models.GameResult {
	res := models.GameResult{
		GameYear:    pa[0].GameYear[0],
		AtBatNumber: abNum,
		Inning:      inning,
		InningTopBot: half,
		Outs:        outs,
		On1b:        bases[0],
		On2b:        bases[1],
		On3b:        bases[2],
		AwayScore:   awayScore,
		HomeScore:   homeScore,
	}
	res.PitchData = append(res.PitchData, pa...)
	return res
}

