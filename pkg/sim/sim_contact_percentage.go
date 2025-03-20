package sim

import (
	"github.com/logananthony/go-baseball/pkg/models"
  "github.com/logananthony/go-baseball/pkg/utils"
)     


// FetchPitcherFrequencies queries and prints pitch data for a pitcher ID
func SimulateContactPercentage(in []models.BatterContactPercentage, stand, pThrows, pitchType string, plateX, plateZ float64) string  {
    
    zone_num := utils.GetPitchZone(plateX, plateZ)
    
    pitch_types := []string{"swinging_strike", "foul", "ball_in_play"} 
    contact_prob := []float64{}
    
    for _, each := range in {
      if each.Stand == stand && each.PThrows == pThrows && each.PitchType == pitchType && each.Zone == zone_num {
        contact_prob = append(contact_prob, each.PctSwingingStrike, each.PctFoul, each.PctBallInPlay)
      }

    }

    sample := utils.WeightedSample(pitch_types, contact_prob)


  return sample

}




