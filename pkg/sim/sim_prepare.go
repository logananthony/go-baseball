package sim

import (
	"database/sql"
	"fmt"

	"github.com/logananthony/go-baseball/pkg/fetcher"
	"github.com/logananthony/go-baseball/pkg/models"
)

// PrepareSimData fetches and builds all setup data for a simulation job.
// Returns error if anything critical is missing.
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
	[]models.BullpenRoleProb, // bullpenRoleProbs
	error,
) {
	// Fetch core game data
	gamePkData, err := fetcher.FetchGameDataByGamePk(db, gameData.GamePk)
	if err != nil || len(gamePkData) == 0 {
		return models.SimData{}, nil, nil, nil, nil, nil, nil, fmt.Errorf("no game data for gamePk %d", gameData.GamePk)
	}
	first := gamePkData[0]
	homeTeam := first.HomeTeamAbbr
	awayTeam := first.AwayTeamAbbr
	season := first.Season
	homeStartingPitcher := first.HomePitcherId
	awayStartingPitcher := first.AwayPitcherId

	// Park Factors
	var parkFactors []models.ParkFactors
	if pf, ok := fetcher.FetchParkFactors(db, homeTeam); ok {
		parkFactors = append(parkFactors, pf)
	} else {
		return models.SimData{}, nil, nil, nil, nil, nil, nil, fmt.Errorf("no park factors for homeTeam %s", homeTeam)
	}

	// PitchingSubProbs, Bullpens
	pitchingSubProbs, err := fetcher.FetchPitchingSubstitutionProbs(db)
	if err != nil {
		return models.SimData{}, nil, nil, nil, nil, nil, nil, fmt.Errorf("could not fetch pitching substitution probs: %v", err)
	}
	homeBullpen := fetcher.FetchBullpenOrder(db, homeTeam, season)
	awayBullpen := fetcher.FetchBullpenOrder(db, awayTeam, season)
	if homeBullpen == nil {
		return models.SimData{}, nil, nil, nil, nil, nil, nil, fmt.Errorf("missing home bullpen for %s, %d", homeTeam, season)
	}
	if awayBullpen == nil {
		return models.SimData{}, nil, nil, nil, nil, nil, nil, fmt.Errorf("missing away bullpen for %s, %d", awayTeam, season)
	}

	// Starter info
	awayStarterInfoSlice, err := fetcher.FetchPlayerInfo(db, awayStartingPitcher)
	if err != nil || len(awayStarterInfoSlice) == 0 {
		return models.SimData{}, nil, nil, nil, nil, nil, nil, fmt.Errorf("no info for away starting pitcher %d", awayStartingPitcher)
	}
	homeStarterInfoSlice, err := fetcher.FetchPlayerInfo(db, homeStartingPitcher)
	if err != nil || len(homeStarterInfoSlice) == 0 {
		return models.SimData{}, nil, nil, nil, nil, nil, nil, fmt.Errorf("no info for home starting pitcher %d", homeStartingPitcher)
	}

	// Lineups
	allGameData, err := fetcher.FetchGameDataByGamePk(db, gameData.GamePk)
	if err != nil || len(allGameData) == 0 {
		return models.SimData{}, nil, nil, nil, nil, nil, nil, fmt.Errorf("no game data for gamePk %d (again)", gameData.GamePk)
	}
	var homeLineup, awayLineup []models.GameDataGamePk
	for _, entry := range allGameData {
		switch entry.Team {
		case "home":
			homeLineup = append(homeLineup, entry)
		case "away":
			awayLineup = append(awayLineup, entry)
		}
	}
	if len(homeLineup) < 9 {
		return models.SimData{}, nil, nil, nil, nil, nil, nil, fmt.Errorf("home lineup too short (%d players)", len(homeLineup))
	}
	if len(awayLineup) < 9 {
		return models.SimData{}, nil, nil, nil, nil, nil, nil, fmt.Errorf("away lineup too short (%d players)", len(awayLineup))
	}

	// Build SimData
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
		MLBParkFactors:      parkFactors,
		HomeTeam:            homeTeam,
	}

	// Add starter info
	simData.PlayerInfo = append(simData.PlayerInfo, awayStarterInfoSlice...)
	simData.PlayerInfo = append(simData.PlayerInfo, homeStarterInfoSlice...)

	// Individual batter/player info for away lineup
	for _, awayBatter := range awayLineup {
		if awayBatter.BattingOrder > 0 && awayBatter.BattingOrder <= 9 {
			awayBatterGameYear := awayBatter.Season
			awayBatterInfo, _ := fetcher.FetchPlayerInfo(db, awayBatter.PlayerId)
			awayBatterSwingProbs, _ := fetcher.FetchBatterSwingPercentage(db, awayBatter.PlayerId, awayBatterGameYear)
			awayBatterContactProbs, _ := fetcher.FetchBatterContactPercentage(db, awayBatter.PlayerId, awayBatterGameYear)
			awayBatterHitProbs, _ := fetcher.FetchBatterHitType(db, awayBatter.PlayerId, awayBatterGameYear)

			simData.PlayerInfo = append(simData.PlayerInfo, awayBatterInfo...)
			simData.BatterSwing = append(simData.BatterSwing, awayBatterSwingProbs...)
			simData.BatterContact = append(simData.BatterContact, awayBatterContactProbs...)
			simData.BatterHitType = append(simData.BatterHitType, awayBatterHitProbs...)
		}
	}

	// Individual batter/player info for home lineup
	for _, homeBatter := range homeLineup {
		if homeBatter.BattingOrder > 0 && homeBatter.BattingOrder <= 9 {
			homeBatterGameYear := homeBatter.Season
			homeBatterInfo, _ := fetcher.FetchPlayerInfo(db, homeBatter.PlayerId)
			homeBatterSwingProbs, _ := fetcher.FetchBatterSwingPercentage(db, homeBatter.PlayerId, homeBatterGameYear)
			homeBatterContactProbs, _ := fetcher.FetchBatterContactPercentage(db, homeBatter.PlayerId, homeBatterGameYear)
			homeBatterHitProbs, _ := fetcher.FetchBatterHitType(db, homeBatter.PlayerId, homeBatterGameYear)

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
		allPitcherYears[i] = season // Could use per-pitcher season if needed
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

	bullpenRoleProbs, _ := fetcher.FetchAllBullpenRoleProbs(db)

	// spew.Dump(bullpenRoleProbs)

	return simData, homeLineup, awayLineup, pitchingSubProbs, homeBullpen, awayBullpen, bullpenRoleProbs, nil
}
