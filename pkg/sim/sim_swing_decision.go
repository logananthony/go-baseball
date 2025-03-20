package sim

import (
	"github.com/logananthony/go-baseball/pkg/models"
  "github.com/logananthony/go-baseball/pkg/utils"
)     


// FetchPitcherFrequencies queries and prints pitch data for a pitcher ID
func SimulateSwingDecision(in []models.BatterSwingPercentage, stand, pThrows, pitchType string, plateX, plateZ float64) bool  {
    
    zone_num := utils.GetPitchZone(plateX, plateZ)
    
    swing_prob := []float64{}
    for _, each := range in {
      if each.Stand == stand && each.PThrows == pThrows && each.PitchType == pitchType && each.Zone == zone_num {
        swing_prob = append(swing_prob, each.SwingPercentage)
      }

    }

    sample := utils.IsSuccess(&swing_prob[0])

  return sample

}



