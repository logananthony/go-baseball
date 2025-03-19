package sim

import (
	"github.com/logananthony/go-baseball/pkg/models"
  "github.com/logananthony/go-baseball/pkg/utils"
)     


// FetchPitcherFrequencies queries and prints pitch data for a pitcher ID
func SimulatePitchType(in []models.PitcherCountPitchFreq, balls, strikes int) string  {

    pitch_types := []string{}
    for _, each := range in {
      if each.BALLS == balls && each.STRIKES == strikes {
        pitch_types = append(pitch_types, each.PITCH_TYPE)
      }

    }

    pitch_weights := []float64{}
    for _, each := range in {
      if each.BALLS == balls && each.STRIKES == strikes {
        pitch_weights = append(pitch_weights, each.FREQUENCY)
      }

    }

  sample := utils.WeightedSample(pitch_types, pitch_weights)

  return sample

}


