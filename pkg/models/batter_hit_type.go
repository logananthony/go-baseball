package models

type BatterHitType struct {
	GameYear       int
	Batter         int
	Stand          string
	PitchType      string
	Zone           int
	PThrows        string
	VelocityBucket string
	Double         float64
	HomeRun        float64
	Out            float64
	Single         float64
	Triple         float64
	N              int
	Level          string
}

