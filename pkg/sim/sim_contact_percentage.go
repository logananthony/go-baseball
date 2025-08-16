package sim

import (
	"fmt"

	"github.com/logananthony/go-baseball/pkg/models"
	"github.com/logananthony/go-baseball/pkg/utils"
)

func SimulateContactPercentage(player map[string]models.BatterContactPercentage, league []models.BatterContactPercentageLeague, stand, pThrows, pitchType string, plateX, plateZ float64) string {

	zoneNum := utils.GetPitchZone(plateX, plateZ)

	pitch_types := []string{"swinging_strike", "foul", "ball_in_play"}
	player_contact_prob := []float64{}
	league_contact_prob := []float64{}
	var contactResult string

	key := fmt.Sprintf("%s|%s|%d|%s", stand, pThrows, zoneNum, pitchType)
	if data, ok := player[key]; ok && data.TotalSwings >= 15 {
		player_contact_prob = []float64{data.PctSwingingStrike, data.PctFoul, data.PctBallInPlay}
	}

	// League-level fallback
	for _, each := range league {
		if each.Stand == stand && each.PThrows == pThrows && each.PitchType == pitchType && each.Zone == zoneNum {
			league_contact_prob = []float64{each.PctSwingingStrike, each.PctFoul, each.PctBallInPlay}
			break
		}
	}

	if len(player_contact_prob) == len(pitch_types) {
		contactResult = utils.WeightedSample(pitch_types, player_contact_prob)
	} else if len(league_contact_prob) == len(pitch_types) {
		contactResult = utils.WeightedSample(pitch_types, league_contact_prob)
	} else {
		contactResult = "foul"
	}

	return contactResult
}
