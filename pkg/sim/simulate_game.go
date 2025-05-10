package sim

import (
	"fmt"

	"github.com/logananthony/go-baseball/pkg/config"
	//"github.com/logananthony/go-baseball/pkg/utils"
	//"github.com/logananthony/go-baseball/pkg/fetcher"
	"github.com/logananthony/go-baseball/pkg/models"
	"github.com/logananthony/go-baseball/pkg/poster"

	//"github.com/logananthony/go-baseball/pkg/utils"
	// "github.com/davecgh/go-spew/spew"
	//"github.com/logananthony/go-baseball/pkg/fetcher"
	"database/sql"

	"github.com/google/uuid"
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

	homePitcher := in[0].HomeStartingPitcherId
  homePitcherGameYear := in[0].HomeStartingPitcherGameYear
	awayPitcher := in[0].AwayStartingPitcherId
  awayPitcherGameYear := in[0].AwayStartingPitcherGameYear

	inning := 7
	awayScore := 3
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

      fmt.Println("Batter #:", awayBatterNumber, "| Event:", awayPaResult[0].EventType[0], " | EV:", awayPaResult[0].ExitVelocity[0], " | LA:", awayPaResult[0].LaunchAngle[0],
      " | SA:", awayPaResult[0].SprayAngle[0],  "| Base State:", awayBaseState[0], awayBaseState[1], awayBaseState[2], "| Score:", awayScore, "-", homeScore)

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

      fmt.Println("Batter #:", homeBatterNumber, "| Event:", homePaResult[0].EventType[0], "| EV:", homePaResult[0].ExitVelocity[0], "| LA:", homePaResult[0].LaunchAngle[0],
      "| SA:", homePaResult[0].SprayAngle[0],  "| Base State:", homeBaseState[0], homeBaseState[1], homeBaseState[2], "| Score:", awayScore, "-", homeScore)

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
    // Bases loaded situation -> force in a run
    if baseState[0] && baseState[1] && baseState[2] {
        baseState[3] = true // Runner scores
        baseState[2] = true // Runner from 2B to 3B
        baseState[1] = true // Runner from 1B to 2B
        baseState[0] = true // Batter to 1B
    } else {
        // Move runners backwards: 2B to 3B, 1B to 2B
        if baseState[1] && baseState[2] {
            baseState[2] = true // Runner from 2B to 3B
            baseState[1] = true // Runner from 1B to 2B
            baseState[0] = true // Batter to 1B
        } else if baseState[1] {
            baseState[2] = baseState[1]
            baseState[1] = baseState[0]
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
          //n := utils.NTrue(baseState[0], baseState[1], baseState[2])
          //if n < 1 {
             // n = 1
          //} 

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


