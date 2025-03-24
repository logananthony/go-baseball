package sim

import(

"github.com/logananthony/go-baseball/pkg/models"
//"reflect"
"fmt"

  
)

func SimulateGame(in []models.GameData) int {

 
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


        for inning < 9 {

              fmt.Println("Inning: ", inning)

              topOuts := 0
              botOuts := 0
              awayBatterNumber := 0
              homeBatterNumber := 0



              for topOuts < 3 {

                  awayBatter := awayLineup[awayBatterNumber % 9]
                  awayPaResult := SimulateAtBat([]models.AtBatData{
                    {
                        GameYear: 2024,
                        PitcherId: homePitcher,
                        BatterId: awayBatter,
                        Strikes: 0,
                        Balls: 0,
                    },
                  })

//                  if awayBatterNumber == 8 {
//                      awayBatterNumber = 0
//                  }
//                  
                  awayBatterNumber++

                  if len(awayPaResult) > 0 && len(awayPaResult[0].HitType) > 0 {
                    if (awayPaResult[0].HitType[0] == "out") || (awayPaResult[0].HitType[0] == "strikeout") {
                          topOuts++
                    }
                  }
                  fmt.Println("Top Inning: ", awayBatterNumber, awayPaResult[0].HitType[0])

              }


              for botOuts < 3 {

                  homeBatter := homeLineup[homeBatterNumber % 9]
                  homePaResult := SimulateAtBat([]models.AtBatData{
                    {
                        GameYear: 2024,
                        PitcherId: awayPitcher,
                        BatterId: homeBatter,
                        Strikes: 0,
                        Balls: 0,
                    },
                  })

//                  if homeBatterNumber == 8 {
//                      homeBatterNumber = 0
//                  }
//                  
                  homeBatterNumber++
                  
                  if len(homePaResult) > 0 && len(homePaResult[0].HitType) > 0 {
                    if (homePaResult[0].HitType[0] == "out") || (homePaResult[0].HitType[0] == "strikeout") {
                          botOuts++
                    }
                  }
                  
                  fmt.Println("Bottom Inning: ", homeBatterNumber, homePaResult[0].HitType[0])

              }

        inning++


      }
        return inning


}

