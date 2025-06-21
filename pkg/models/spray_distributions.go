package models

import "database/sql"

type SprayDistribution struct {
	Batter            int
	GameYear          int
	Stand             string
	PThrows           string
	Outcome           *string
	Zone              *int
	EVBucket          *string
	LaunchAngleBucket *string
	Skew              sql.NullFloat64
	Mean              sql.NullFloat64
	Std               sql.NullFloat64
	N                 int
	Level             string
}
