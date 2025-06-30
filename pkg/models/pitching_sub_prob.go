package models

type PitchingSubstitutionProb struct {
	Inning           int
	RunsScoredGame   int
	RunsScoredInning int
	Role             int
	PullProbability  float64
}
