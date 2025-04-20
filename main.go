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
    HomeBatter1Id: 542932, HomeBatter1GameYear: 2024,           
    HomeBatter2Id: 665828, HomeBatter2GameYear: 2024,     
    HomeBatter3Id: 592450, HomeBatter3GameYear: 2024,         
    HomeBatter4Id: 665862, HomeBatter4GameYear: 2024,   
    HomeBatter5Id: 605204, HomeBatter5GameYear: 2024,           
    HomeBatter6Id: 691176, HomeBatter6GameYear: 2024,     
    HomeBatter7Id: 701538, HomeBatter7GameYear: 2024,           
    HomeBatter8Id: 663330, HomeBatter8GameYear: 2024,
    HomeBatter9Id: 663757, HomeBatter9GameYear: 2024,

    HomeStartingPitcherId: 543037, HomeStartingPitcherGameYear: 2024,

    // Dodgers (Away)
    AwayBatter1Id: 605141, AwayBatter1GameYear: 2024,         
    AwayBatter2Id: 669242, AwayBatter2GameYear: 2024,
    AwayBatter3Id: 518692, AwayBatter3GameYear: 2024,      
    AwayBatter4Id: 606192, AwayBatter4GameYear: 2024,
    AwayBatter5Id: 571771, AwayBatter5GameYear: 2024,
    AwayBatter6Id: 595281, AwayBatter6GameYear: 2024,
    AwayBatter7Id: 666158, AwayBatter7GameYear: 2024,
    AwayBatter8Id: 676439, AwayBatter8GameYear: 2024,
    AwayBatter9Id: 605131, AwayBatter9GameYear: 2024,        

    AwayStartingPitcherId: 607192, AwayStartingPitcherGameYear: 2024,
},

})



      _=game_result


    //fmt.Println(game_result)

     //spew.Dump(game_result)





  }

