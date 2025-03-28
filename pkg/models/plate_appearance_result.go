package models


type PlateAppearanceResult struct { 
  PitcherId []int
  BatterId []int
  BatterStands []string
  PitcherThrows []string
  Strikes []int 
  Balls []int
  PitchCount []int
  PitchType []string
  PlateX []float64
  PlateZ []float64
  Zone []int
  Velocity []float64
  IsStrike []bool
  IsSwing []bool
  IsContact []string
  EventType []string
}




