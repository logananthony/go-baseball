package models

type PitchingSubstitutionProb struct {
	Inning           int
	RunsScoredGame   int
	RunsScoredInning int
	TotalAppearances int
	TotalPulled      int
	PullProbability  float64
}

