package sim

import (
	"fmt"

	"github.com/logananthony/go-baseball/pkg/models"
	"github.com/logananthony/go-baseball/pkg/utils"
)

// // SimulateBatterHitType selects the best matching batted ball outcome for a batter
// // using a specificity scoring system across Stand, PThrows, PitchType, Zone, and VelocityBucket.
// func SimulateBatterHitType(in []models.BatterHitType, stand, pThrows, pitchType string, plateX, plateZ, velocity float64, teamAbbr string, pf []models.ParkFactors) string {
// 	zone_num := utils.GetPitchZone(plateX, plateZ)
// 	velo_bucket := utils.GetVelocityBucket(velocity)

// 	var bestScore int
// 	var selected *models.BatterHitType

// 	for i := range in {
// 		score := 0
// 		if in[i].Stand == stand {
// 			score++
// 		}
// 		if in[i].PThrows == pThrows {
// 			score++
// 		}
// 		if in[i].PitchType == pitchType {
// 			score++
// 		}
// 		if in[i].Zone == zone_num {
// 			score++
// 		}
// 		if in[i].VelocityBucket == velo_bucket {
// 			score++
// 		}

// 		if score > bestScore {
// 			bestScore = score
// 			selected = &in[i]
// 		}
// 	}

// 	if selected == nil || bestScore < 2 { // Optional: require minimum specificity
// 		return "out"
// 	}

// 	outcomes := []string{"single", "double", "triple", "home_run", "out"}
// 	probs := []float64{
// 		selected.Single,
// 		selected.Double,
// 		selected.Triple,
// 		selected.HomeRun,
// 		selected.Out,
// 	}

// 	return utils.WeightedSample(outcomes, probs)
// }

func SimulateBatterHitType(in []models.BatterHitType, stand, pThrows, pitchType string, plateX, plateZ, velocity float64, teamAbbr string, pf []models.ParkFactors) string {
	zone_num := utils.GetPitchZone(plateX, plateZ)
	velo_bucket := utils.GetVelocityBucket(velocity)

	var bestScore int
	var selected *models.BatterHitType

	for i := range in {
		score := 0
		if in[i].Stand == stand {
			score++
		}
		if in[i].PThrows == pThrows {
			score++
		}
		if in[i].PitchType == pitchType {
			score++
		}
		if in[i].Zone == zone_num {
			score++
		}
		if in[i].VelocityBucket == velo_bucket {
			score++
		}

		if score > bestScore {
			bestScore = score
			selected = &in[i]
		}
	}

	if selected == nil || bestScore < 2 {
		return "out"
	}

	// ─── Park Factor Lookup ───────────────────────────────────────────────
	var park models.ParkFactors
	found := false
	for _, p := range pf {
		if p.Team == teamAbbr {
			park = p
			found = true
			break
		}
	}

	// ─── Probabilities ────────────────────────────────────────────────────
	probs := []float64{
		selected.Single,
		selected.Double,
		selected.Triple,
		selected.HomeRun,
		selected.Out,
	}

	outcomes := []string{"single", "double", "triple", "home_run", "out"}

	// ─── Park Factor Adjustments ──────────────────────────────────────────
	if found {
		fmt.Printf("🌍 Applying park factors for %s (%s): %+v\n", teamAbbr, park.Venue, park)

		before := make([]float64, len(probs))
		copy(before, probs)

		probs[0] *= float64(park.OneB) / 100.0
		probs[1] *= float64(park.TwoB) / 100.0
		probs[2] *= float64(park.ThreeB) / 100.0
		probs[3] *= float64(park.HR) / 100.0
		probs[4] *= 2.0 - float64(park.H)/100.0

		for i := range outcomes {
			fmt.Printf("🔢 %s: raw=%.4f → adjusted=%.4f\n", outcomes[i], before[i], probs[i])
		}
	} else {
		fmt.Printf("⚠️  No park factor found for team: %s\n", teamAbbr)
	}

	return utils.WeightedSample(outcomes, probs)
}
