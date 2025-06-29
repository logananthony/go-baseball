// package sim

// import (
// 	//"fmt"
// 	"github.com/logananthony/go-baseball/pkg/models"
// 	"github.com/logananthony/go-baseball/pkg/utils"
// )

// // SimulateBatterHitType safely handles NULLable fields when matching
// func SimulateBatterHitType(in []models.BatterHitType, stand, pThrows, pitchType string, plateX, plateZ, velocity float64) string {
// 	zone_num := utils.GetPitchZone(plateX, plateZ)
// 	velo_bucket := utils.GetVelocityBucket(velocity)

// 	var selected *models.BatterHitType

//   for _, each := range in {
//     if each.Stand == stand ||
//        each.PThrows == pThrows ||
//        each.PitchType == pitchType ||
//        each.Zone == zone_num ||
//        each.VelocityBucket == velo_bucket {
//       selected = &each
//       break
//     }
//   }

// 	if selected == nil {
// 		return "out"
// 	}

// 	// Directly use the raw probabilities from the selected match
// 	outcomes := []string{"single", "double", "triple", "home_run", "out"}

//   probs := []float64{
// 	selected.Single,
// 	selected.Double,
// 	selected.Triple,
// 	selected.HomeRun,
// 	selected.Out,
// }

// 	//fmt.Println("Raw probabilities:", probs)
// 	return utils.WeightedSample(outcomes, probs)
// }

package sim

import (
	"github.com/logananthony/go-baseball/pkg/models"
	"github.com/logananthony/go-baseball/pkg/utils"
)

// SimulateBatterHitType selects the best matching batted ball outcome for a batter
// using a specificity scoring system across Stand, PThrows, PitchType, Zone, and VelocityBucket.
func SimulateBatterHitType(in []models.BatterHitType, stand, pThrows, pitchType string, plateX, plateZ, velocity float64) string {
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

	if selected == nil || bestScore < 2 { // Optional: require minimum specificity
		return "out"
	}

	outcomes := []string{"single", "double", "triple", "home_run", "out"}
	probs := []float64{
		selected.Single,
		selected.Double,
		selected.Triple,
		selected.HomeRun,
		selected.Out,
	}

	// fmt.Printf("🎯 HitType Weights: single=%.2f double=%.2f triple=%.2f HR=%.2f out=%.2f\n",
	// 	selected.Single, selected.Double, selected.Triple, selected.HomeRun, selected.Out)

	return utils.WeightedSample(outcomes, probs)
}
