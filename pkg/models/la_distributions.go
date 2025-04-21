package models

type LADistribution struct {
	GameYear      int
	Batter        int
	Stand         string
	PThrows       string
	Outcome       *string
	Zone          *int
	EVBucket      *string
	Skew          float64
	Mean          float64
	Std           float64
	N             int
	Level         string
}

