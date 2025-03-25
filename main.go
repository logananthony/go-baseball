package main

import (
    //"github.com/logananthony/go-baseball/pkg/fetcher"
    "github.com/logananthony/go-baseball/pkg/models"
    "github.com/logananthony/go-baseball/pkg/sim"
    //"github.com/logananthony/go-baseball/pkg/fetcher"
    //"github.com/logananthony/go-baseball/pkg/config"
    //"github.com/davecgh/go-spew/spew"
    "fmt"
    //"encoding/json"
)

func main() {

//  sim_results := sim.SimulateAtBat([]models.AtBatData{
//    {
//        GameYear: 2024,
//        PitcherId: 628317,
//        BatterId: 665742,
//        Strikes: 0,
//        Balls: 0,
//    },
//  })




  //spew.Dump(sim_results)


  game_result := sim.SimulateGame([]models.GameData{
              {
                HomeBatter1: 665742,
                HomeBatter2: 665742,
                HomeBatter3: 665742,
                HomeBatter4: 665742,
                HomeBatter5: 665742,
                HomeBatter6: 665742,
                HomeBatter7: 665742,
                HomeBatter8: 665742,
                HomeBatter9: 665742,

                AwayBatter1: 665742,
                AwayBatter2: 665742,
                AwayBatter3: 665742,
                AwayBatter4: 665742,
                AwayBatter5: 665742,
                AwayBatter6: 665742,
                AwayBatter7: 665742,
                AwayBatter8: 665742,
                AwayBatter9: 665742,

                HomeStartingPitcher: 571945, 
                AwayStartingPitcher: 571945,


              },
            })



    fmt.Println(game_result)



  }

