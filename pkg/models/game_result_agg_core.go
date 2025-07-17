package models

import "time"

type GameResultAggCore struct {
	GamePk           int64   // bigint
	TotalOver35      float64 // double precision
	TotalOver45      float64
	TotalOver65      float64
	TotalOver75      float64
	TotalOver85      float64
	TotalOver95      float64
	TotalOver105     float64
	SpreadMinus25    float64
	SpreadMinus15    float64
	SpreadPlus15     float64
	SpreadPlus25     float64
	MoneylineHomeWin float64
	MoneylineAwayWin float64
	GameDate         time.Time // timestamp without time zone
	HomeTeamAbbr     string
	AwayTeamAbbr     string
	HomePitcherName  string
	AwayPitcherName  string
	TotalOver15      float64
	TotalOver25      float64
	TotalOver55      float64
	TotalOver115     float64
	TotalOver125     float64
	SpreadMinus55    float64
	SpreadMinus45    float64
	SpreadMinus35    float64
	SpreadPlus35     float64
	SpreadPlus45     float64
	SpreadPlus55     float64
	HomeTotalOver05  float64 // numeric
	HomeTotalOver15  float64
	HomeTotalOver25  float64
	HomeTotalOver35  float64
	HomeTotalOver45  float64
	HomeTotalOver55  float64
	HomeTotalOver65  float64
	AwayTotalOver05  float64
	AwayTotalOver15  float64
	AwayTotalOver25  float64
	AwayTotalOver35  float64
	AwayTotalOver45  float64
	AwayTotalOver55  float64
	AwayTotalOver65  float64
	StdTotalRuns     float64 // numeric
	StdHomeScore     float64
	StdAwayScore     float64
	MlVar            float64
	StdSpread        float64
}
