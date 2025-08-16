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
	MLBParkFactors      []ParkFactors
	HomeTeam            string

	// Precomputed lookups for faster simulation
	PlayerInfoMap       map[int]MLBPlayerInfo
	PitcherCovMap       map[int][]PitcherCovarianceMean
	PitcherPitchFreqMap map[int]map[string][]PitcherCountPitchFreq
	BatterSwingMap      map[int]map[string]BatterSwingPercentage
	BatterContactMap    map[int]map[string]BatterContactPercentage
	BatterHitTypeMap    map[int]map[string]BatterHitType
}
