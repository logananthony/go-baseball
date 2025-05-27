package sim

import (
	"github.com/logananthony/go-baseball/pkg/models"
  "github.com/logananthony/go-baseball/pkg/utils"
)     


func SimulateSwingDecision(player []models.BatterSwingPercentage, league []models.BatterSwingPercentageLeague, stand, pThrows, pitchType string, plateX, plateZ float64) bool {
	
  zoneNum := utils.GetPitchZone(plateX, plateZ)

	var playerSwing *float64 = nil
	for _, each := range player {
		if each.Stand == stand && each.PThrows == pThrows && each.PitchType == pitchType && each.Zone == zoneNum {
			if each.TotalPitches >= 25 {
				playerSwing = &each.SwingPercentage
				break
			}
		}
	}

	if playerSwing != nil {
		return utils.IsSuccess(playerSwing)
	}

	for _, each := range league {
		if each.Stand == stand && each.PThrows == pThrows && each.PitchType == pitchType && each.Zone == zoneNum {
			val := float64(each.SwingPercentage)
			return utils.IsSuccess(&val)
		}
	}

	return false
}

