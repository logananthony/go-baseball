package sim

import (
	"fmt"

	"github.com/logananthony/go-baseball/pkg/models"
	"github.com/logananthony/go-baseball/pkg/utils"
)

func SimulateSwingDecision(player map[string]models.BatterSwingPercentage, league []models.BatterSwingPercentageLeague, stand, pThrows, pitchType string, plateX, plateZ float64) bool {

	zoneNum := utils.GetPitchZone(plateX, plateZ)

	key := fmt.Sprintf("%s|%s|%d|%s", stand, pThrows, zoneNum, pitchType)
	if data, ok := player[key]; ok && data.TotalPitches >= 25 {
		return utils.IsSuccess(&data.SwingPercentage)
	}

	for _, each := range league {
		if each.Stand == stand && each.PThrows == pThrows && each.PitchType == pitchType && each.Zone == zoneNum {
			val := float64(each.SwingPercentage)
			return utils.IsSuccess(&val)
		}
	}

	return false
}
