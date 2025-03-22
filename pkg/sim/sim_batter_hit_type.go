package sim

import (
	//"fmt"
	"github.com/logananthony/go-baseball/pkg/models"
	"github.com/logananthony/go-baseball/pkg/utils"
)

// SimulateBatterHitType safely handles NULLable fields when matching
func SimulateBatterHitType(in []models.BatterHitType, stand, pThrows, pitchType string, plateX, plateZ, velocity float64) string {
	zone_num := utils.GetPitchZone(plateX, plateZ)
	velo_bucket := utils.GetVelocityBucket(velocity)

	var selected *models.BatterHitType

  for _, each := range in {
    if each.Stand == stand || 
       each.PThrows == pThrows || 
       each.PitchType == pitchType || 
       each.Zone == zone_num || 
       each.VelocityBucket == velo_bucket {
      selected = &each
      break
    }
  }


	if selected == nil {
		return "out"
	}

	// Directly use the raw probabilities from the selected match
	outcomes := []string{"single", "double", "triple", "home_run", "out"}

  probs := []float64{
	selected.Single,
	selected.Double,
	selected.Triple,
	selected.HomeRun,
	selected.Out,
}

	//fmt.Println("Raw probabilities:", probs)
	return utils.WeightedSample(outcomes, probs)
}

