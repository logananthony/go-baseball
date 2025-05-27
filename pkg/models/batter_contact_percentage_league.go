package models

type BatterContactPercentageLeague struct {
  Stand             string    `csv:"stand"`
  PThrows           string    `csv:"p_throws"`
  Zone              int       `csv:"zone"`
  PitchType         string    `csv:"pitch_type"`
  BallInPlay        int       `csv:"ball_in_play"`
  Foul              int       `csv:"foul"`
  SwingingStrike    int       `csv:"swinging_strike"`
  TotalSwings       int       `csv:"total_swings"`
  PctSwingingStrike float64   `csv:"pct_swinging_strike"`
  PctFoul           float64   `csv:"pct_foul"`
  PctBallInPlay     float64   `csv:"pct_ball_in_play"`
} 

