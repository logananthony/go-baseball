package models

type BatterSwingPercentage struct {
	GameYear        int
	Batter          int
	Stand           string
	PThrows         string
	Zone            int
	PitchType       string
	TotalPitches    int
	TotalSwings     int
	SwingPercentage float64
}
