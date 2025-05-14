package main
 
import (
    "github.com/logananthony/go-baseball/pkg/models"
    "github.com/logananthony/go-baseball/pkg/sim"
    "github.com/joho/godotenv"
    "log"
    //"os"
)

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found")
    }
}

func main() {

  err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, relying on system environment variables")
    }

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

        HomeStartingPitcherId: 678394, HomeStartingPitcherGameYear: 2024, // Gerrit Cole
        HomeBullpen1Id: 695243, HomeBullpen1GameYear: 2024,
        HomeBullpen2Id: 695243, HomeBullpen2GameYear: 2024,
        HomeBullpen3Id: 695243, HomeBullpen3GameYear: 2024,
        HomeBullpen4Id: 695243, HomeBullpen4GameYear: 2024,
        HomeBullpen5Id: 695243, HomeBullpen5GameYear: 2024,
        HomeBullpen6Id: 695243, HomeBullpen6GameYear: 2024,
        HomeBullpen7Id: 695243, HomeBullpen7GameYear: 2024,
        HomeBullpen8Id: 695243, HomeBullpen8GameYear: 2024,
        HomeBullpen9Id: 695243, HomeBullpen9GameYear: 2024,

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

        AwayStartingPitcherId: 678394, AwayStartingPitcherGameYear: 2024, // Walker Buehler
        AwayBullpen1Id: 695243, AwayBullpen1GameYear: 2024,
        AwayBullpen2Id: 695243, AwayBullpen2GameYear: 2024,
        AwayBullpen3Id: 695243, AwayBullpen3GameYear: 2024,
        AwayBullpen4Id: 695243, AwayBullpen4GameYear: 2024,
        AwayBullpen5Id: 695243, AwayBullpen5GameYear: 2024,
        AwayBullpen6Id: 695243, AwayBullpen6GameYear: 2024,
        AwayBullpen7Id: 695243, AwayBullpen7GameYear: 2024,
        AwayBullpen8Id: 695243, AwayBullpen8GameYear: 2024,
        AwayBullpen9Id: 695243, AwayBullpen9GameYear: 2024,
                
    },
})

      _=game_result


  }

