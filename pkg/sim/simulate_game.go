package sim

import(

"github.com/logananthony/go-baseball/pkg/models"
//"reflect"
"fmt"

  
)

func SimulateGame(in []models.GameData) []models.GameResult {

 
        homeLineup := []int{
                       in[0].HomeBatter1,
                       in[0].HomeBatter2,
                       in[0].HomeBatter3,
                       in[0].HomeBatter4,
                       in[0].HomeBatter5,
                       in[0].HomeBatter6,
                       in[0].HomeBatter7,
                       in[0].HomeBatter8,
                       in[0].HomeBatter9,
                      }


        awayLineup := []int{
                       in[0].AwayBatter1,
                       in[0].AwayBatter2,
                       in[0].AwayBatter3,
                       in[0].AwayBatter4,
                       in[0].AwayBatter5,
                       in[0].AwayBatter6,
                       in[0].AwayBatter7,
                       in[0].AwayBatter8,
                       in[0].AwayBatter9,
                      }

        homePitcher := in[0].HomeStartingPitcher
        awayPitcher := in[0].AwayStartingPitcher


        inning := 1
        awayScore := 0
        homeScore := 0
        awayBatterNumber := 0
        homeBatterNumber := 0
        atBatNumber := 0

        //var plateAppRes []models.GameResult

        var gameRes []models.GameResult
        //var paRes PlateAppearanceResult

        for inning <= 1 {

              topOuts := 0
              botOuts := 0
              awayBaseState := []bool {false, false, false, false}
              homeBaseState := []bool {false, false, false, false}

              fmt.Println("Top Inning: ", inning)

              for topOuts < 3 {


                  awayBatterNumber = awayBatterNumber % 9
                  awayBatter := awayLineup[awayBatterNumber]
                  awayPaResult := SimulateAtBat([]models.PlateAppearanceData{
                    {
                        GameYear: 2024,
                        PitcherId: homePitcher,
                        BatterId: awayBatter,
                        Strikes: 0,
                        Balls: 0,
                    },
                  })

                  atBatNumber += 1


                 if len(awayPaResult) > 0 && len(awayPaResult[0].EventType) > 0 {
                      switch awayPaResult[0].EventType[0] {
                      case "walk":
                          if awayBaseState[0] || (awayBaseState[0] && awayBaseState[2]) {
                                  awayBaseState[0] = awayBaseState[1]
                          } else if awayBaseState[0] && awayBaseState[1] {
                                  awayBaseState[0] = awayBaseState[1]
                                  awayBaseState[1] = awayBaseState[2]
                          } else if awayBaseState[0] && awayBaseState[1] && awayBaseState[2] {
                                  awayBaseState[0] = awayBaseState[1]
                                  awayBaseState[1] = awayBaseState[2]
                                  awayBaseState[2] = awayBaseState[3]
                          } else if !awayBaseState[0] {
                                  awayBaseState[0] = true      
                          }
                      case "single":
                          awayBaseState[3] = awayBaseState[2] 
                          awayBaseState[2] = awayBaseState[1]
                          awayBaseState[1] = awayBaseState[0] 
                          awayBaseState[0] = true         

                      case "double":
                          awayBaseState[3] = awayBaseState[1] 
                          awayBaseState[2] = awayBaseState[0] 
                          awayBaseState[1] = true
                          awayBaseState[0] = false

                      case "triple":
                          awayBaseState[3] = awayBaseState[0] || awayBaseState[1] || awayBaseState[2] 
                          awayBaseState[2] = true
                          awayBaseState[1] = false
                          awayBaseState[0] = false

                      case "home_run":
                          runs := 1 
                          for i := 0; i <= 2; i++ {
                              if awayBaseState[i] {
                                  runs++
                                  awayBaseState[i] = false 
                              }
                          }
                          awayScore += runs
                      }

                      if awayBaseState[3] {
                          awayScore++
                          awayBaseState[3] = false

                      }
                  }

                  gameRes = append(gameRes, models.GameResult{
                      GameYear: awayPaResult[0].GameYear[0],
                      AwayScore: awayScore, 
                      HomeScore: homeScore,
                      AtBatNumber: atBatNumber,
                      Inning: inning,
                      InningTopBot: "Top",
                      Outs: topOuts,
                      On1b: awayBaseState[0],
                      On2b: awayBaseState[1],
                      On3b: awayBaseState[2],
                  })
                  gameRes[0].PitchData = append(gameRes[0].PitchData, awayPaResult...)

 

                  awayBatterNumber++


                  if len(awayPaResult) > 0 && len(awayPaResult[0].EventType) > 0 {
                    if (awayPaResult[0].EventType[0] == "out") || (awayPaResult[0].EventType[0] == "strikeout") {
                          topOuts++
                    }
                  }

                  fmt.Println("Batter #: ", awayBatterNumber, " | Event: ",  awayPaResult[0].EventType[0], " | Base State: ", 
                              awayBaseState[0], awayBaseState[1], awayBaseState[2],   " | Score: ", awayScore, "-", homeScore)

              }


              fmt.Println("Bottom Inning: ", inning)


              for botOuts < 3 {

                  homeBatterNumber = homeBatterNumber % 9
                  homeBatter := homeLineup[homeBatterNumber]
                  homePaResult := SimulateAtBat([]models.PlateAppearanceData{
                    {
                        GameYear: 2024,
                        PitcherId: awayPitcher,
                        BatterId: homeBatter,
                        Strikes: 0,
                        Balls: 0,
                    },
                  })

                  atBatNumber += 1


                 if len(homePaResult) > 0 && len(homePaResult[0].EventType) > 0 {
                      switch homePaResult[0].EventType[0] {
                      case "walk":
                          if homeBaseState[0] || (homeBaseState[0] && homeBaseState[2]) {
                                  homeBaseState[0] = homeBaseState[1]
                          } else if homeBaseState[0] && homeBaseState[1] {
                                  homeBaseState[0] = homeBaseState[1]
                                  homeBaseState[1] = homeBaseState[2]
                          } else if homeBaseState[0] && homeBaseState[1] && homeBaseState[2] {
                                  homeBaseState[0] = homeBaseState[1]
                                  homeBaseState[1] = homeBaseState[2]
                                  homeBaseState[2] = homeBaseState[3]
                          } else if !homeBaseState[0] {
                                  homeBaseState[0] = true      
                          }
                      case "single":
                          homeBaseState[3] = homeBaseState[2] 
                          homeBaseState[2] = homeBaseState[1]
                          homeBaseState[1] = homeBaseState[0] 
                          homeBaseState[0] = true      
                      case "double":
                          homeBaseState[3] = homeBaseState[1] 
                          homeBaseState[2] = homeBaseState[0] 
                          homeBaseState[1] = true
                          homeBaseState[0] = false
                      case "triple":
                          homeBaseState[3] = homeBaseState[0] || homeBaseState[1] || homeBaseState[2] 
                          homeBaseState[2] = true
                          homeBaseState[1] = false
                          homeBaseState[0] = false
                      case "home_run":
                          runs := 1 
                          for i := 0; i <= 2; i++ {
                              if homeBaseState[i] {
                                  runs++
                                  homeBaseState[i] = false 
                              }
                          }
                          homeScore += runs
                      }


                      if homeBaseState[3] {
                          homeScore++
                          homeBaseState[3] = false 
                      }
                  }

                  gameRes = append(gameRes, models.GameResult{
                      GameYear: homePaResult[0].GameYear[0],
                      AwayScore: awayScore, 
                      HomeScore: homeScore, 
                      AtBatNumber: atBatNumber,
                      Inning: inning,
                      InningTopBot: "Bot",
                      Outs: botOuts,
                      On1b: homeBaseState[0],
                      On2b: homeBaseState[1],
                      On3b: homeBaseState[2],
                  })
                  gameRes[0].PitchData = append(gameRes[0].PitchData, homePaResult...)


                
                  homeBatterNumber++
                  
                  if len(homePaResult) > 0 && len(homePaResult[0].EventType) > 0 {
                    if (homePaResult[0].EventType[0] == "out") || (homePaResult[0].EventType[0] == "strikeout") {
                          botOuts++
                    }
                  }

                  fmt.Println("Batter #: ", homeBatterNumber, " | Event: ",  homePaResult[0].EventType[0], " | Base State: ", 
                              homeBaseState[0], homeBaseState[1], homeBaseState[2],   " | Score: ", awayScore, "-", homeScore)

              }

        inning++


      }
        return gameRes


}

