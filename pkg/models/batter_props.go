package models

import "time"

type BatterProps struct {
	BatterID     int64
	GamePk       int64
	BatterName   string
	BattingOrder int
	Team         string
	GameDate     time.Time

	NumSimulations int

	// Hits (all hit types combined)
	ProbOver05Hits    float64
	ProbOver15Hits    float64
	Over05HitsLower95 float64
	Over05HitsUpper95 float64
	Over15HitsLower95 float64
	Over15HitsUpper95 float64
	AvgHits           float64
	IqrHits           float64
	Q80Hits           float64

	// Singles
	ProbOver05Singles    float64
	ProbOver15Singles    float64
	Over05SinglesLower95 float64
	Over05SinglesUpper95 float64
	Over15SinglesLower95 float64
	Over15SinglesUpper95 float64
	AvgSingles           float64
	IqrSingles           float64
	Q80Singles           float64

	// Doubles
	ProbOver05Doubles    float64
	ProbOver15Doubles    float64
	Over05DoublesLower95 float64
	Over05DoublesUpper95 float64
	Over15DoublesLower95 float64
	Over15DoublesUpper95 float64
	AvgDoubles           float64
	IqrDoubles           float64
	Q80Doubles           float64

	// Triples
	ProbOver05Triples    float64
	ProbOver15Triples    float64
	Over05TriplesLower95 float64
	Over05TriplesUpper95 float64
	Over15TriplesLower95 float64
	Over15TriplesUpper95 float64
	AvgTriples           float64
	IqrTriples           float64
	Q80Triples           float64

	// Home runs
	ProbOver05Homeruns    float64
	ProbOver15Homeruns    float64
	Over05HomerunsLower95 float64
	Over05HomerunsUpper95 float64
	Over15HomerunsLower95 float64
	Over15HomerunsUpper95 float64
	AvgHomeruns           float64
	IqrHomeruns           float64
	Q80Homeruns           float64
	BattingAvg            float64
	SluggingPct           float64
	OnBasePct             float64
	StrikeoutPct          float64
	WalkPct               float64

	// Runs scored
	ProbOver05Runs    float64
	ProbOver15Runs    float64
	Over05RunsLower95 float64
	Over05RunsUpper95 float64
	Over15RunsLower95 float64
	Over15RunsUpper95 float64
	AvgRuns           float64
	IqrRuns           float64
	Q80Runs           float64

	// RBIs
	ProbOver05RBI    float64
	ProbOver15RBI    float64
	Over05RBILower95 float64
	Over05RBIUpper95 float64
	Over15RBILower95 float64
	Over15RBIUpper95 float64
	AvgRBI           float64
	IqrRBI           float64
	Q80RBI           float64
}
