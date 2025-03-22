package models
import (
   // "database/sql"
    //"fmt"
    _ "github.com/lib/pq"
)

type BatterHitType struct {
    GameYear        int
    Batter          int
    Stand           string
    PThrows         string
    Zone            int
    PitchType       string
    VelocityBucket  string
    Double          float64
    HomeRun         float64
    Out             float64
    Single          float64
    Triple          float64
    N               int
    Level           string
}


