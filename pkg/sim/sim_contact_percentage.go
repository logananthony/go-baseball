package sim

import (
	"github.com/logananthony/go-baseball/pkg/models"
	"github.com/logananthony/go-baseball/pkg/utils"
	// "fmt"
)

func SimulateContactPercentage(player []models.BatterContactPercentage, league []models.BatterContactPercentageLeague, stand, pThrows, pitchType string, plateX, plateZ float64) string {

	zoneNum := utils.GetPitchZone(plateX, plateZ)

	pitch_types := []string{"swinging_strike", "foul", "ball_in_play"}
	player_contact_prob := []float64{}
	league_contact_prob := []float64{}
	var contactResult string

	// Player-level check with more lenient sample size for 2025
	for _, each := range player {
		if each.Stand == stand && each.PThrows == pThrows && each.PitchType == pitchType && each.Zone == zoneNum {
			if each.TotalSwings >= 10 { // Reduced from 25 for 2025 data
				player_contact_prob = []float64{each.PctSwingingStrike, each.PctFoul, each.PctBallInPlay}
				// fmt.Println("✅ Using player-level data")
				break
			}
		}
	}

	// League-level fallback
	for _, each := range league {
		if each.Stand == stand && each.PThrows == pThrows && each.PitchType == pitchType && each.Zone == zoneNum {
			league_contact_prob = []float64{each.PctSwingingStrike, each.PctFoul, each.PctBallInPlay}
			// fmt.Println("📉 Using league-level fallback")
			break
		}
	}

	// Try more general league data if specific not found
	if len(league_contact_prob) == 0 {
		for _, each := range league {
			if each.Stand == stand && each.PThrows == pThrows && each.PitchType == pitchType {
				league_contact_prob = []float64{each.PctSwingingStrike, each.PctFoul, each.PctBallInPlay}
				break
			}
		}
	}

	if len(player_contact_prob) == len(pitch_types) {
		contactResult = utils.WeightedSample(pitch_types, player_contact_prob)
	} else if len(league_contact_prob) == len(pitch_types) {
		contactResult = utils.WeightedSample(pitch_types, league_contact_prob)
	} else {
		// fmt.Println("🟥 No data available, defaulting to 'ball_in_play'")
		contactResult = "foul"
	}

	return contactResult
}
