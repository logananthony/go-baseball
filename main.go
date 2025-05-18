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

    // Testing real game (2024 Tigers @ Red Sox - 746954)
    {   
        HomeBatter1Id: 663837, HomeBatter1GameYear: 2024, // Matt Vierling
        HomeBatter2Id: 682985, HomeBatter2GameYear: 2024, // Riley Greene
        HomeBatter3Id: 592192, HomeBatter3GameYear: 2024, // Mark Canha
        HomeBatter4Id: 570482, HomeBatter4GameYear: 2024, // Gio Urshela
        HomeBatter5Id: 690993, HomeBatter5GameYear: 2024, // Colt Keith
        HomeBatter6Id: 668731, HomeBatter6GameYear: 2024, // Akil Badoo
        HomeBatter7Id: 679529, HomeBatter7GameYear: 2024, // Spencer Torkelson
        HomeBatter8Id: 595879, HomeBatter8GameYear: 2024, // Javy Baez
        HomeBatter9Id: 608348, HomeBatter9GameYear: 2024, // Carson Kelly

        HomeStartingPitcherId: 656427, HomeStartingPitcherGameYear: 2024, // Jack Flaherty
        HomeBullpen1Id: 571946, HomeBullpen1GameYear: 2024, // Shelby Miller
        HomeBullpen2Id: 689225, HomeBullpen2GameYear: 2024, // Beau Brieske
        HomeBullpen3Id: 676684, HomeBullpen3GameYear: 2024, // Will Vest
        HomeBullpen4Id: 663947, HomeBullpen4GameYear: 2024, // Tyler Holton
        HomeBullpen5Id: 656412, HomeBullpen5GameYear: 2024, // Alex Faedo
        HomeBullpen6Id: 605177, HomeBullpen6GameYear: 2024, // Andrew Chafin
        HomeBullpen7Id: 669724, HomeBullpen7GameYear: 2024, // Brenan Hanifee
        HomeBullpen8Id: 666159, HomeBullpen8GameYear: 2024, // Matt Manning
        HomeBullpen9Id: 680744, HomeBullpen9GameYear: 2024, // Ty Madden

        AwayBatter1Id: 680776, AwayBatter1GameYear: 2024, // Jarren Duran 
        AwayBatter2Id: 677800, AwayBatter2GameYear: 2024, // Wilyer Abreu
        AwayBatter3Id: 608701, AwayBatter3GameYear: 2024, // Rob Refsnyder
        AwayBatter4Id: 646240, AwayBatter4GameYear: 2024, // Rafael Devers
        AwayBatter5Id: 657136, AwayBatter5GameYear: 2024, // Connor Wong
        AwayBatter6Id: 643265, AwayBatter6GameYear: 2024, // Garrett Cooper
        AwayBatter7Id: 624512, AwayBatter7GameYear: 2024, // Reese McGuire
        AwayBatter8Id: 687093, AwayBatter8GameYear: 2024, // Vaugh Grissom
        AwayBatter9Id: 666152, AwayBatter9GameYear: 2024, // David Hamilton

        AwayStartingPitcherId: 601713, AwayStartingPitcherGameYear: 2024, // Nick Pivetta
        AwayBullpen1Id: 502624, AwayBullpen1GameYear: 2024, // Chase Anderson
        AwayBullpen2Id: 686580, AwayBullpen2GameYear: 2024, // Justin Slaten
        AwayBullpen3Id: 670174, AwayBullpen3GameYear: 2024, // Josh Winckowski
        AwayBullpen4Id: 455119, AwayBullpen4GameYear: 2024, // Chris Martin
        AwayBullpen5Id: 657514, AwayBullpen5GameYear: 2024, // Brennan Bernardino
        AwayBullpen6Id: 677161, AwayBullpen6GameYear: 2024, // Zack Kelly
        AwayBullpen7Id: 669711, AwayBullpen7GameYear: 2024, // Greg Weissert
        AwayBullpen8Id: 445276, AwayBullpen8GameYear: 2024, // Kenley Jansen
        AwayBullpen9Id: 592155, AwayBullpen9GameYear: 2024, // Cam Booser
                
    },
})

      _=game_result


  }

