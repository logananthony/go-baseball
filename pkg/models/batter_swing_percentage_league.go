package models

type BatterSwingPercentageLeague struct {
  Stand           string  `csv:"stand"`
  PThrows         string  `csv:"p_throws"`
  Zone            int     `csv:"zone"`
  PitchType       string  `csv:"pitch_type"`
  TotalPitches    int     `csv:"total_pitches"`
  TotalSwings     int     `csv:"total_swings"`
  SwingPercentage float64 `csv:"swing_percentage"`

} 
