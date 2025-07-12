package sim

import (
	"database/sql"

	"github.com/logananthony/go-baseball/pkg/fetcher"
	"github.com/logananthony/go-baseball/pkg/models"
)

// PrepareSimData fetches and builds all setup data for a simulation job.
// It does NOT put pitchingSubProbs, homeBullpen, or awayBullpen inside SimData struct;
// these are returned as separate values and should be passed as function args to SimulateGame.
func PrepareSimData(
	db *sql.DB,
	gameData models.GameData,
) (
	models.SimData, // simData
	[]models.GameDataGamePk, // homeLineup
	[]models.GameDataGamePk, // awayLineup
	[]models.PitchingSubstitutionProb, // pitchingSubProbs
	*models.BullpenOrder, // homeBullpen
	*models.BullpenOrder, // awayBullpen
	[]models.BullpenRoleProb, // bullpenRoleProbs (NEW)
	error,
) {
	gamePkData, _ := fetcher.FetchGameDataByGamePk(db, gameData.GamePk)
	first := gamePkData[0]
	homeTeam := first.HomeTeamAbbr
	awayTeam := first.AwayTeamAbbr
	season := first.Season
	homeStartingPitcher := first.HomePitcherId
	awayStartingPitcher := first.AwayPitcherId

	// Build SimData (the struct that will be reused across all simulations for this job)
	simData := models.SimData{
		PlayerInfo:          []models.MLBPlayerInfo{},
		LeagueSwing:         fetcher.FetchBatterSwingPercentageLeague(),
		LeagueContact:       fetcher.FetchBatterContactPercentageLeague(),
		LeaguePitchCovMeans: fetcher.FetchPitcherCovarianceMeanLeague(),
		BatterSwing:         []models.BatterSwingPercentage{},
		BatterContact:       []models.BatterContactPercentage{},
		BatterHitType:       []models.BatterHitType{},
		PitcherPitchFreq:    []models.PitcherCountPitchFreq{},
		PitcherCovMeans:     []models.PitcherCovarianceMean{},
		BatterEVDist:        []models.EVDistribution{},
		BatterLADist:        []models.LADistribution{},
		BatterSprayDist:     []models.SprayDistribution{},
		MLBParkFactors:      []models.ParkFactors{},
		HomeTeam:            homeTeam,
	}

	// Park Factors
	if pf, ok := fetcher.FetchParkFactors(db, homeTeam); ok {
		simData.MLBParkFactors = append(simData.MLBParkFactors, pf)
	}

	// ---- Fetch the extra setup objects (not in SimData) ----
	pitchingSubProbs, _ := fetcher.FetchPitchingSubstitutionProbs(db)
	homeBullpen := fetcher.FetchBullpenOrder(db, homeTeam, season)
	awayBullpen := fetcher.FetchBullpenOrder(db, awayTeam, season)
	// --------------------------------------------------------

	// Starter info
	awayStarterInfoSlice, _ := fetcher.FetchPlayerInfo(db, awayStartingPitcher)
	homeStarterInfoSlice, _ := fetcher.FetchPlayerInfo(db, homeStartingPitcher)
	simData.PlayerInfo = append(simData.PlayerInfo, awayStarterInfoSlice...)
	simData.PlayerInfo = append(simData.PlayerInfo, homeStarterInfoSlice...)

	// Lineups (all game data, split into home/away)
	allGameData, _ := fetcher.FetchGameDataByGamePk(db, gameData.GamePk)
	var homeLineup, awayLineup []models.GameDataGamePk
	for _, entry := range allGameData {
		switch entry.Team {
		case "home":
			homeLineup = append(homeLineup, entry)
		case "away":
			awayLineup = append(awayLineup, entry)
		}
	}

	// Add individual batter/player info for away lineup
	for _, awayBatter := range awayLineup {
		if awayBatter.BattingOrder > 0 && awayBatter.BattingOrder <= 9 {
			awayBatterGameYear := awayBatter.Season
			awayBatterInfo, _ := fetcher.FetchPlayerInfo(db, awayBatter.PlayerId)
			awayBatterSwingProbs, _ := fetcher.FetchBatterSwingPercentage(db, awayBatter.PlayerId, awayBatterGameYear)
			awayBatterContactProbs, _ := fetcher.FetchBatterContactPercentage(db, awayBatter.PlayerId, awayBatterGameYear)
			awayBatterHitProbs, _ := fetcher.FetchBatterHitType(db, awayBatter.PlayerId, awayBatterGameYear)
			// You can add additional fetches here for EV/LA/Spray if needed

			simData.PlayerInfo = append(simData.PlayerInfo, awayBatterInfo...)
			simData.BatterSwing = append(simData.BatterSwing, awayBatterSwingProbs...)
			simData.BatterContact = append(simData.BatterContact, awayBatterContactProbs...)
			simData.BatterHitType = append(simData.BatterHitType, awayBatterHitProbs...)
		}
	}

	// Add individual batter/player info for home lineup
	for _, homeBatter := range homeLineup {
		if homeBatter.BattingOrder > 0 && homeBatter.BattingOrder <= 9 {
			homeBatterGameYear := homeBatter.Season
			homeBatterInfo, _ := fetcher.FetchPlayerInfo(db, homeBatter.PlayerId)
			homeBatterSwingProbs, _ := fetcher.FetchBatterSwingPercentage(db, homeBatter.PlayerId, homeBatterGameYear)
			homeBatterContactProbs, _ := fetcher.FetchBatterContactPercentage(db, homeBatter.PlayerId, homeBatterGameYear)
			homeBatterHitProbs, _ := fetcher.FetchBatterHitType(db, homeBatter.PlayerId, homeBatterGameYear)
			// You can add additional fetches here for EV/LA/Spray if needed

			simData.PlayerInfo = append(simData.PlayerInfo, homeBatterInfo...)
			simData.BatterSwing = append(simData.BatterSwing, homeBatterSwingProbs...)
			simData.BatterContact = append(simData.BatterContact, homeBatterContactProbs...)
			simData.BatterHitType = append(simData.BatterHitType, homeBatterHitProbs...)
		}
	}

	// Build pitcher lists for each team (starters + bullpen)
	homePitchers := []int{
		homeStartingPitcher,
		homeBullpen.PlayerID1,
		homeBullpen.PlayerID2,
		homeBullpen.PlayerID3,
		homeBullpen.PlayerID4,
		homeBullpen.PlayerID5,
		homeBullpen.PlayerID6,
		homeBullpen.PlayerID7,
		homeBullpen.PlayerID8,
	}
	awayPitchers := []int{
		awayStartingPitcher,
		awayBullpen.PlayerID1,
		awayBullpen.PlayerID2,
		awayBullpen.PlayerID3,
		awayBullpen.PlayerID4,
		awayBullpen.PlayerID5,
		awayBullpen.PlayerID6,
		awayBullpen.PlayerID7,
		awayBullpen.PlayerID8,
	}
	allPitchers := append(homePitchers, awayPitchers...)
	allPitcherYears := make([]int, len(allPitchers))
	for i := range allPitchers {
		allPitcherYears[i] = season // You could use per-pitcher season if needed
	}

	for i, pitcherId := range allPitchers {
		year := allPitcherYears[i]
		pitchFreqsR := fetcher.FetchPitcherFrequencies(db, pitcherId, "R", year)
		pitchFreqsL := fetcher.FetchPitcherFrequencies(db, pitcherId, "L", year)
		simData.PitcherPitchFreq = append(simData.PitcherPitchFreq, pitchFreqsR...)
		simData.PitcherPitchFreq = append(simData.PitcherPitchFreq, pitchFreqsL...)
		pitcherCov := fetcher.FetchPitcherCovarianceMean(db, int64(pitcherId), int64(year))
		simData.PitcherCovMeans = append(simData.PitcherCovMeans, pitcherCov...)
		pitcherInfo, _ := fetcher.FetchPlayerInfo(db, pitcherId)
		simData.PlayerInfo = append(simData.PlayerInfo, pitcherInfo...)
	}

	bullpenRoleProbs, _ := fetcher.FetchBullpenRoleProbs(db, homeTeam, awayTeam, season)

	return simData, homeLineup, awayLineup, pitchingSubProbs, homeBullpen, awayBullpen, bullpenRoleProbs, nil

}
