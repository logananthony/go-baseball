package models

import "time"

type GameResultAggCore struct {
	GamePk int64

	// --- Game Totals (point + CI) ---
	TotalOver15         float64
	TotalOver15Lower95  float64
	TotalOver15Upper95  float64
	TotalOver25         float64
	TotalOver25Lower95  float64
	TotalOver25Upper95  float64
	TotalOver35         float64
	TotalOver35Lower95  float64
	TotalOver35Upper95  float64
	TotalOver45         float64
	TotalOver45Lower95  float64
	TotalOver45Upper95  float64
	TotalOver55         float64
	TotalOver55Lower95  float64
	TotalOver55Upper95  float64
	TotalOver65         float64
	TotalOver65Lower95  float64
	TotalOver65Upper95  float64
	TotalOver75         float64
	TotalOver75Lower95  float64
	TotalOver75Upper95  float64
	TotalOver85         float64
	TotalOver85Lower95  float64
	TotalOver85Upper95  float64
	TotalOver95         float64
	TotalOver95Lower95  float64
	TotalOver95Upper95  float64
	TotalOver105        float64
	TotalOver105Lower95 float64
	TotalOver105Upper95 float64
	TotalOver115        float64
	TotalOver115Lower95 float64
	TotalOver115Upper95 float64
	TotalOver125        float64
	TotalOver125Lower95 float64
	TotalOver125Upper95 float64

	// --- Home Team Totals (point + CI) ---
	HomeTotalOver05        float64
	HomeTotalOver05Lower95 float64
	HomeTotalOver05Upper95 float64
	HomeTotalOver15        float64
	HomeTotalOver15Lower95 float64
	HomeTotalOver15Upper95 float64
	HomeTotalOver25        float64
	HomeTotalOver25Lower95 float64
	HomeTotalOver25Upper95 float64
	HomeTotalOver35        float64
	HomeTotalOver35Lower95 float64
	HomeTotalOver35Upper95 float64
	HomeTotalOver45        float64
	HomeTotalOver45Lower95 float64
	HomeTotalOver45Upper95 float64
	HomeTotalOver55        float64
	HomeTotalOver55Lower95 float64
	HomeTotalOver55Upper95 float64
	HomeTotalOver65        float64
	HomeTotalOver65Lower95 float64
	HomeTotalOver65Upper95 float64

	// --- Away Team Totals (point + CI) ---
	AwayTotalOver05        float64
	AwayTotalOver05Lower95 float64
	AwayTotalOver05Upper95 float64
	AwayTotalOver15        float64
	AwayTotalOver15Lower95 float64
	AwayTotalOver15Upper95 float64
	AwayTotalOver25        float64
	AwayTotalOver25Lower95 float64
	AwayTotalOver25Upper95 float64
	AwayTotalOver35        float64
	AwayTotalOver35Lower95 float64
	AwayTotalOver35Upper95 float64
	AwayTotalOver45        float64
	AwayTotalOver45Lower95 float64
	AwayTotalOver45Upper95 float64
	AwayTotalOver55        float64
	AwayTotalOver55Lower95 float64
	AwayTotalOver55Upper95 float64
	AwayTotalOver65        float64
	AwayTotalOver65Lower95 float64
	AwayTotalOver65Upper95 float64

	// --- Spread Lines (point + CI) ---
	SpreadMinus55        float64
	SpreadMinus55Lower95 float64
	SpreadMinus55Upper95 float64
	SpreadMinus45        float64
	SpreadMinus45Lower95 float64
	SpreadMinus45Upper95 float64
	SpreadMinus35        float64
	SpreadMinus35Lower95 float64
	SpreadMinus35Upper95 float64
	SpreadMinus25        float64
	SpreadMinus25Lower95 float64
	SpreadMinus25Upper95 float64
	SpreadMinus15        float64
	SpreadMinus15Lower95 float64
	SpreadMinus15Upper95 float64
	SpreadPlus15         float64
	SpreadPlus15Lower95  float64
	SpreadPlus15Upper95  float64
	SpreadPlus25         float64
	SpreadPlus25Lower95  float64
	SpreadPlus25Upper95  float64
	SpreadPlus35         float64
	SpreadPlus35Lower95  float64
	SpreadPlus35Upper95  float64
	SpreadPlus45         float64
	SpreadPlus45Lower95  float64
	SpreadPlus45Upper95  float64
	SpreadPlus55         float64
	SpreadPlus55Lower95  float64
	SpreadPlus55Upper95  float64

	// --- Moneylines and their confidence intervals ---
	MoneylineHomeWin float64
	MlHomeWinLower95 float64
	MlHomeWinUpper95 float64
	MoneylineAwayWin float64

	// --- Meta/game info ---
	GameDate        time.Time
	HomeTeamAbbr    string
	AwayTeamAbbr    string
	HomePitcherName string
	AwayPitcherName string

	// --- Distribution stats ---
	StdTotalRuns float64
	IqrTotalRuns float64
	Q80TotalRuns float64

	StdHomeScore     float64
	IqrHomeScore     float64
	Q80HomeScore     float64
	HomeScoreLower95 float64
	HomeScoreUpper95 float64

	StdAwayScore     float64
	IqrAwayScore     float64
	Q80AwayScore     float64
	AwayScoreLower95 float64
	AwayScoreUpper95 float64

	MlVar float64

	StdSpread     float64
	IqrSpread     float64
	Q80Spread     float64
	SpreadLower95 float64
	SpreadUpper95 float64

	AvgTotalRuns float64
	AvgHomeScore float64
	AvgAwayScore float64
}
