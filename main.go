package main
 
import (
    //"github.com/logananthony/go-baseball/pkg/fetcher"
    "github.com/logananthony/go-baseball/pkg/models"
    "github.com/logananthony/go-baseball/pkg/sim"
    //"github.com/logananthony/go-baseball/pkg/fetcher"
    //"github.com/logananthony/go-baseball/pkg/config"
    "github.com/davecgh/go-spew/spew"
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
        HomeBatter1Id: 542932, HomeBatter1FullName: "Jon Berti", HomeBatter1GameYear: 2024,           
        HomeBatter2Id: 665828, HomeBatter2FullName: "Oswaldo Cabrera", HomeBatter2GameYear: 2024,     
        HomeBatter3Id: 592450, HomeBatter3FullName: "Aaron Judge", HomeBatter3GameYear: 2024,         
        HomeBatter4Id: 665862, HomeBatter4FullName: "Jazz Chisholm Jr.", HomeBatter4GameYear: 2024,   
        HomeBatter5Id: 605204, HomeBatter5FullName: "J.D. Davis", HomeBatter5GameYear: 2024,           
        HomeBatter6Id: 691176, HomeBatter6FullName: "Jasson Domínguez", HomeBatter6GameYear: 2024,     
        HomeBatter7Id: 668752, HomeBatter7FullName: "Duke Ellis", HomeBatter7GameYear: 2024,           
        HomeBatter8Id: 663330, HomeBatter8FullName: "Jahmai Jones", HomeBatter8GameYear: 2024,
        HomeBatter9Id: 663757, HomeBatter9FullName: "Trent Grisham", HomeBatter9GameYear: 2024,

        HomeStartingPitcherId: 677076, HomeStartingPitcherFullName: "Clayton Andrews", HomeStartingPitcherGameYear: 2024,

        // Dodgers (Away)
        AwayBatter1Id: 605141, AwayBatter1FullName: "Mookie Betts", AwayBatter1GameYear: 2024,         
        AwayBatter2Id: 669242, AwayBatter2FullName: "Tommy Edman", AwayBatter2GameYear: 2024,
        AwayBatter3Id: 518692, AwayBatter3FullName: "Freddie Freeman", AwayBatter3GameYear: 2024,      
        AwayBatter4Id: 606192, AwayBatter4FullName: "Teoscar Hernández", AwayBatter4GameYear: 2024,
        AwayBatter5Id: 571771, AwayBatter5FullName: "Enrique Hernández", AwayBatter5GameYear: 2024,
        AwayBatter6Id: 595281, AwayBatter6FullName: "Kevin Kiermaier", AwayBatter6GameYear: 2024,
        AwayBatter7Id: 666158, AwayBatter7FullName: "Gavin Lux", AwayBatter7GameYear: 2024,
        AwayBatter8Id: 676439, AwayBatter8FullName: "Hunter Feduccia", AwayBatter8GameYear: 2024,
        AwayBatter9Id: 605131, AwayBatter9FullName: "Austin Barnes", AwayBatter9GameYear: 2024,        

        AwayStartingPitcherId: 607455, AwayStartingPitcherFullName: "Anthony Banda", AwayStartingPitcherGameYear: 2024,
    },
})



      _=game_result


    //fmt.Println(game_result)

     spew.Dump(game_result)





  }

