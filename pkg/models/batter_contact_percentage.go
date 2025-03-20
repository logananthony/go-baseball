package models

type BatterContactPercentage struct {
    GameYear            int 
    Batter              int 
    Stand               string  
    PThrows             string
    Zone                int  
    PitchType           string  
    SwingingStrike      int 
    Foul                int  
    BallInPlay          int  
    TotalSwings         int 
    PctSwingingStrike   float64 
    PctFoul             float64
    PctBallInPlay       float64 
}

