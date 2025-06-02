package models

type SimData struct {
	PlayerInfo          []MLBPlayerInfo
	LeagueSwing         []BatterSwingPercentageLeague
	LeagueContact       []BatterContactPercentageLeague
	LeaguePitchCovMeans []PitcherCovarianceMeanLeague
	BatterSwing         []BatterSwingPercentage
	BatterContact       []BatterContactPercentage
	BatterHitType       []BatterHitType
	PitcherPitchFreq    []PitcherCountPitchFreq
	PitcherCovMeans     []PitcherCovarianceMean
	BatterEVDist        []EVDistribution
	BatterLADist        []LADistribution
	BatterSprayDist     []SprayDistribution
}
