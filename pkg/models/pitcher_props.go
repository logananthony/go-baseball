package models

import "time"

type PitcherProps struct {
	PitcherID      int64
	GamePk         int64
	PitcherName    string
	GameDate       time.Time
	NumSimulations int

	// Strikeout lines
	ProbOver25K  float64
	ProbOver35K  float64
	ProbOver45K  float64
	ProbOver55K  float64
	ProbOver65K  float64
	ProbOver75K  float64
	ProbOver85K  float64
	ProbOver95K  float64
	ProbOver105K float64
	ProbOver115K float64
	ProbOver125K float64

	Over25KLower95  float64
	Over25KUpper95  float64
	Over35KLower95  float64
	Over35KUpper95  float64
	Over45KLower95  float64
	Over45KUpper95  float64
	Over55KLower95  float64
	Over55KUpper95  float64
	Over65KLower95  float64
	Over65KUpper95  float64
	Over75KLower95  float64
	Over75KUpper95  float64
	Over85KLower95  float64
	Over85KUpper95  float64
	Over95KLower95  float64
	Over95KUpper95  float64
	Over105KLower95 float64
	Over105KUpper95 float64
	Over115KLower95 float64
	Over115KUpper95 float64
	Over125KLower95 float64
	Over125KUpper95 float64

	AvgStrikeouts float64
	IqrStrikeouts float64
	Q80Strikeouts float64
}
