package models


type GameResult struct {

    //GameYear int
    AwayScore int
    HomeScore int
    AtBatNumber int
    Inning int
    InningTopBot string
    Outs int
    On1b bool
    On2b bool 
    On3b bool
    PitchData []PlateAppearanceResult

    
  
}

