package utils

import (
	"github.com/logananthony/go-baseball/pkg/models"
)

// FilterPitcherExitDistributions filters the exit probability data
func FilterPitcherExitDistributions(
	data []models.PitchingSubstitutionProb,
	pitcherRole string,
	teamDeficit int,
	pitchBin string,
) (float64, bool) {
	for _, row := range data {
		if row.PitcherRole == pitcherRole &&
			row.PitchBin == pitchBin {

			if row.PitcherRole == "starter" {
				return DeficitBooster(teamDeficit, row.PctOfRole), true
				// return row.CumulativePctOfRole, true
			}
			return row.CumulativePctOfRole, true
			// fmt.Println("boosted prpb: ", DeficitBooster(teamDeficit, row.PctOfRole))
			// return DeficitBooster(teamDeficit, row.PctOfRole), true
			// return row.PctOfRole, true
		}
	}
	return 0.0, false
}
