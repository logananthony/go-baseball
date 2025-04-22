package main
 
import (
    //"github.com/logananthony/go-baseball/pkg/fetcher"
    "github.com/logananthony/go-baseball/pkg/models"
    "github.com/logananthony/go-baseball/pkg/sim"
    //"github.com/logananthony/go-baseball/pkg/fetcher"
    //"github.com/logananthony/go-baseball/pkg/config"
    //"github.com/davecgh/go-spew/spew"
    //_"github.com/lib/pq"
	  //"database/sql"
    //"github.com/jmoiron/sqlx"
    //"log"
    //"fmt"
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
        // Yankees (Home)
        HomeBatter1Id: 592450, HomeBatter1GameYear: 2024, // Aaron Judge
        HomeBatter2Id: 665862, HomeBatter2GameYear: 2024, // Juan Soto
        HomeBatter3Id: 605204, HomeBatter3GameYear: 2024, // Giancarlo Stanton
        HomeBatter4Id: 543760, HomeBatter4GameYear: 2024, // Anthony Rizzo
        HomeBatter5Id: 547180, HomeBatter5GameYear: 2024, // DJ LeMahieu
        HomeBatter6Id: 666158, HomeBatter6GameYear: 2024, // Gleyber Torres
        HomeBatter7Id: 660271, HomeBatter7GameYear: 2024, // Oswaldo Cabrera
        HomeBatter8Id: 663757, HomeBatter8GameYear: 2024, // Anthony Volpe
        HomeBatter9Id: 701538, HomeBatter9GameYear: 2024, // Austin Wells

        HomeStartingPitcherId: 601713, HomeStartingPitcherGameYear: 2024, // Gerrit Cole

        // Dodgers (Away)
        AwayBatter1Id: 660271, AwayBatter1GameYear: 2024, // Mookie Betts
        AwayBatter2Id: 660271, AwayBatter2GameYear: 2024, // Shohei Ohtani
        AwayBatter3Id: 547180, AwayBatter3GameYear: 2024, // Freddie Freeman
        AwayBatter4Id: 606192, AwayBatter4GameYear: 2024, // Will Smith
        AwayBatter5Id: 605141, AwayBatter5GameYear: 2024, // Max Muncy
        AwayBatter6Id: 669242, AwayBatter6GameYear: 2024, // Gavin Lux
        AwayBatter7Id: 518692, AwayBatter7GameYear: 2024, // Chris Taylor
        AwayBatter8Id: 571771, AwayBatter8GameYear: 2024, // James Outman
        AwayBatter9Id: 595281, AwayBatter9GameYear: 2024, // Miguel Vargas

        AwayStartingPitcherId: 656302, AwayStartingPitcherGameYear: 2024, // Walker Buehler
    },
})



      _=game_result


    //fmt.Println(game_result)

     //spew.Dump(game_result)





  }

