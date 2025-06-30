package sim

import (
	"github.com/logananthony/go-baseball/pkg/models"
	"github.com/logananthony/go-baseball/pkg/utils"
)

func SimulatePitchType(in []models.PitcherCountPitchFreq, balls, strikes int) string {
	pitch_types, pitch_weights := filterByCount(in, balls, strikes)

	// Try fallback counts if no exact match
	if len(pitch_types) == 0 {
		candidates := [][2]int{
			{balls, strikes - 1},
			{balls, strikes + 1},
			{balls - 1, strikes},
			{balls + 1, strikes},
		}
		for _, bsc := range candidates {
			pitch_types, pitch_weights = filterByCount(in, bsc[0], bsc[1])
			if len(pitch_types) > 0 {
				break
			}
		}
	}

	// If still nothing, use all data as final fallback
	if len(pitch_types) == 0 {
		for _, each := range in {
			pitch_types = append(pitch_types, each.PITCH_TYPE)
			pitch_weights = append(pitch_weights, each.FREQUENCY)
		}
	}

	return utils.WeightedSample(pitch_types, pitch_weights)
}

func filterByCount(in []models.PitcherCountPitchFreq, balls, strikes int) ([]string, []float64) {
	var types []string
	var weights []float64
	for _, each := range in {
		if each.BALLS == balls && each.STRIKES == strikes {
			types = append(types, each.PITCH_TYPE)
			weights = append(weights, each.FREQUENCY)
		}
	}
	return types, weights
}
