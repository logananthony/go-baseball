package models


type SprayDistribution struct {
    Batter            int
    GameYear          int
    Stand             string
    PThrows           string
    Outcome           *string
    Zone              *int       
    EVBucket          *string
    LaunchAngleBucket *string
    Skew              float64
    Mean              float64
    Std               float64
    N                 int
    Level             string
}
