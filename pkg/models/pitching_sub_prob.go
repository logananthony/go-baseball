package models

type PitchingSubstitutionProb struct {
	PitcherRole string
	// TeamDeficit  int
	PitchBin            string
	NumPitchers         int
	PctOfRole           float64
	CumulativePctOfRole float64
}
