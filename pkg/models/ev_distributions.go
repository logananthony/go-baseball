package models

import "database/sql"

type EVDistribution struct {
	GameYear       int
	Batter         int
	Stand          string
	PThrows        string
	Outcome        sql.NullString
	PitchType      sql.NullString
	Zone           sql.NullInt32
	VelocityBucket sql.NullString
	Skew           float64
	Mean           float64
	Std            float64
	N              int
	Level          string
}

