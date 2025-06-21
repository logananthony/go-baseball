package models

import "database/sql"

type LADistribution struct {
	GameYear int
	Batter   int
	Stand    string
	PThrows  string
	Outcome  *string
	Zone     *int
	EVBucket *string
	Skew     sql.NullFloat64
	Mean     sql.NullFloat64
	Std      sql.NullFloat64
	N        int
	Level    string
}
