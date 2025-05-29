package sim

import (
	"fmt"
	"math/rand"
	//"github.com/davecgh/go-spew/spew"
	"github.com/logananthony/go-baseball/pkg/config"
	//"github.com/logananthony/go-baseball/pkg/fetcher"
	"database/sql"

	"github.com/google/uuid"
	"github.com/logananthony/go-baseball/pkg/fetcher"
	"github.com/logananthony/go-baseball/pkg/models"
	"github.com/logananthony/go-baseball/pkg/poster"
	"github.com/logananthony/go-baseball/pkg/utils"
	//"log"
)

func SimulateGame(in []models.GameData) []models.GameResult {

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

  inning := 7
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

 
	var gameRes []models.GameResult

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
			awayPaResult := SimulateAtBat([]models.PlateAppearanceData{{
				BatterGameYear: awayBatterGameYear,
				BatterId:  awayBatter,
				PitcherGameYear:  homePitcherGameYear,
				PitcherId: homePitcher,
				Strikes:   0,
				Balls:     0,
			}})

			atBatNumber++
			awayScore, awayBaseState, topOuts = ProcessPlateAppearance(
				awayPaResult, awayScore, awayBaseState, topOuts,
			)

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
                  "| Event:", awayPaResult[0].EventType[0], 
                  "| EV:", awayPaResult[0].ExitVelocity[0], 
                  "| LA:", awayPaResult[0].LaunchAngle[0],
                  "| SA:", awayPaResult[0].SprayAngle[0],  
                  "| Base State:", awayBaseState[0], awayBaseState[1], awayBaseState[2], 
                  "| Score:", awayScore, "-", homeScore)
                  //"| Inning Runs:", inningRunsHome, 
                  //"| Pull Prob Home:", pullProbHome, pitcherPulledHome)

			gameRes = append(gameRes, BuildGameResult(awayPaResult, atBatNumber, inning, "Top", topOuts, awayBaseState, awayScore, homeScore))
			awayBatterNumber++
		}

		if inning >= 9 && homeScore > awayScore {
      postGameResults(gameRes, db)
			fmt.Println("Home team wins:", homeScore, "-", awayScore)
			break
		}

		fmt.Println("Bottom Inning:", inning)

		for botOuts < 3 {
      homeBatterNumber = homeBatterNumber % 9
			homeBatter := homeLineup[homeBatterNumber]
      homeBatterGameYear := homeLineupGameYear[homeBatterNumber]
			homePaResult := SimulateAtBat([]models.PlateAppearanceData{{
				BatterGameYear: homeBatterGameYear,
				BatterId:  homeBatter,
				PitcherGameYear:  awayPitcherGameYear,
				PitcherId: awayPitcher,
				Strikes:   0,
				Balls:     0,
			}})

			atBatNumber++
			homeScore, homeBaseState, botOuts = ProcessPlateAppearance(
				homePaResult, homeScore, homeBaseState, botOuts,
			)

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
                  "| Event:", homePaResult[0].EventType[0], 
                  "| EV:", homePaResult[0].ExitVelocity[0],
                  "| LA:", homePaResult[0].LaunchAngle[0],
                  "| SA:", homePaResult[0].SprayAngle[0],  
                  "| Base State:", homeBaseState[0], homeBaseState[1], homeBaseState[2], 
                  "| Score:", awayScore, "-", homeScore)

			gameRes = append(gameRes, BuildGameResult(homePaResult, atBatNumber, inning, "Bot", botOuts, homeBaseState, awayScore, homeScore))
			homeBatterNumber++

			if inning >= 9 && homeScore > awayScore {
				fmt.Println("Home team wins (walk-off):", homeScore, "-", awayScore)
        postGameResults(gameRes, db)
				return gameRes
			} 
		}

		// If 9 or later and not tied, game ends
		if inning >= 9 && homeScore != awayScore {
			fmt.Println("Away team wins:", awayScore, "-", homeScore)
      postGameResults(gameRes, db)
			break
		}

		inning++
	}

	return gameRes
}

func ProcessPlateAppearance(paResult []models.PlateAppearanceResult, score int, baseState []bool, outs int) (int, []bool, int) {

  db := config.ConnectDB()
  defer db.Close()

	if len(paResult) == 0 || len(paResult[0].EventType) == 0 {
		return score, baseState, outs 
	}

	switch paResult[0].EventType[0] {
  case "walk":
    if baseState[0] && baseState[1] && baseState[2] {
        baseState[3] = true 
        baseState[2] = true 
        baseState[1] = true 
        baseState[0] = true // Batter to B
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
            // No runners at all
            baseState[0] = true
        }
    }

      case "single", "double", "triple":

          priorBaseState := baseState
          newBaseState := [4]bool{false, false, false, false}

          for i := range baseState {

              if priorBaseState[i] {
                  var basesMoved int
                  if paResult[0].EventType[0] == "single" {
                      basesMoved = 1
                  } else if paResult[0].EventType[0] == "double" {
                      basesMoved = 2
                  } else {
                      basesMoved = 3
                  }

                  if i + basesMoved < 3 {
                      newBaseState[i + basesMoved] = true 
                  } else {
                      score++
                  }
                  
              } else {
                      if paResult[0].EventType[0] == "single" {
                          newBaseState[0] = true
                      } else if paResult[0].EventType[0] == "double" {
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

func BuildGameResult(pa []models.PlateAppearanceResult, abNum, inning int, half string, outs int, bases []bool, awayScore, homeScore int) models.GameResult {
	res := models.GameResult{
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


func postGameResults(gameRes []models.GameResult, db *sql.DB) {
  gameId := uuid.New().String()
	gameYear := 2024

	for _, result := range gameRes {
		// Optional debug
		if len(result.PitchData) >= 0 {
			fmt.Printf("DEBUG: Inserting %d pitches for AtBat #%d\n", len(result.PitchData[0].PitchType), result.AtBatNumber)
		}

		err := poster.InsertGameResult(db, gameId, gameYear, result)
		if err != nil {
			fmt.Println("Error inserting game result:", err)
		}
	}
}


