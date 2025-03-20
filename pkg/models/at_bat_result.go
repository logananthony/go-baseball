package models


type AtBatResult struct { 
  GameYear int
  PitcherId int
  BatterId int
  Strikes []int
  Balls []int
  //PitchCount int
  PitchType []string
  PlateX []float64
  PlateZ []float64
  Velocity []float64
  IsStrike []bool
  IsSwing []bool
  IsContact []string
}




