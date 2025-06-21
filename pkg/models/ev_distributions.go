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
	Skew           sql.NullFloat64
	Mean           sql.NullFloat64
	Std            sql.NullFloat64
	N              sql.NullInt32
	Level          string
}
