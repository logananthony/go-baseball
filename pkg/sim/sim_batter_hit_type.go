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

// 	// Fix: Use AND (&&) instead of OR (||) for proper matching
// 	// Also add sample size validation
// 	for _, each := range in {
// 		if each.Stand == stand &&
// 		   each.PThrows == pThrows &&
// 		   each.PitchType == pitchType &&
// 		   each.Zone == zone_num &&
// 		   each.VelocityBucket == velo_bucket {
// 			// Check if we have sufficient sample size
// 			if each.N >= 10 {
// 				selected = &each
// 				break
// 			}
// 		}
// 	}

// 	// If no specific match found, try to find a more general match (without velocity bucket)
// 	if selected == nil {
// 		for _, each := range in {
// 			if each.Stand == stand &&
// 			   each.PThrows == pThrows &&
// 			   each.PitchType == pitchType &&
// 			   each.Zone == zone_num {
// 				if each.N >= 10 {
// 					selected = &each
// 					break
// 				}
// 			}
// 		}
// 	}

// 	// If still no match, try even more general (without pitch type)
// 	if selected == nil {
// 		for _, each := range in {
// 			if each.Stand == stand &&
// 			   each.PThrows == pThrows &&
// 			   each.Zone == zone_num {
// 				if each.N >= 10 {
// 					selected = &each
// 					break
// 				}
// 			}
// 		}
// 	}

// 	// If still no match, try most general (just stand and pitcher throws)
// 	if selected == nil {
// 		for _, each := range in {
// 			if each.Stand == stand &&
// 			   each.PThrows == pThrows {
// 				if each.N >= 10 {
// 					selected = &each
// 					break
// 				}
// 			}
// 		}
// 	}

// 	if selected == nil {
// 		return "out"
// 	}

// 	// Directly use the raw probabilities from the selected match
// 	outcomes := []string{"single", "double", "triple", "home_run", "out"}

// 	probs := []float64{
// 		selected.Single,
// 		selected.Double,
// 		selected.Triple,
// 		selected.HomeRun,
// 		selected.Out,
// 	}

// 	//fmt.Println("Raw probabilities:", probs)
// 	return utils.WeightedSample(outcomes, probs)
// }

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

	// Try to find the most specific match with sufficient sample size
	for _, each := range in {
		if each.Stand == stand &&
			each.PThrows == pThrows &&
			each.PitchType == pitchType &&
			each.Zone == zone_num &&
			each.VelocityBucket == velo_bucket {
			// For 2025 data, we might need to be more lenient with sample sizes
			if each.N >= 5 {
				selected = &each
				break
			}
		}
	}

	// If no specific match found, try without velocity bucket
	if selected == nil {
		for _, each := range in {
			if each.Stand == stand &&
				each.PThrows == pThrows &&
				each.PitchType == pitchType &&
				each.Zone == zone_num {
				if each.N >= 5 {
					selected = &each
					break
				}
			}
		}
	}

	// If still no match, try without pitch type
	if selected == nil {
		for _, each := range in {
			if each.Stand == stand &&
				each.PThrows == pThrows &&
				each.Zone == zone_num {
				if each.N >= 5 {
					selected = &each
					break
				}
			}
		}
	}

	// If still no match, try most general (just stand and pitcher throws)
	if selected == nil {
		for _, each := range in {
			if each.Stand == stand &&
				each.PThrows == pThrows {
				if each.N >= 10 {
					selected = &each
					break
				}
			}
		}
	}

	// If still no match, try any data for this batter (last resort)
	if selected == nil {
		for _, each := range in {
			if each.N >= 1 {
				selected = &each
				break
			}
		}
	}

	if selected == nil {
		// Return a reasonable default distribution
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

	// Ensure probabilities sum to a reasonable value
	total := 0.0
	for _, p := range probs {
		total += p
	}

	// If probabilities are too low or invalid, use league averages
	if total < 0.1 {
		// League average distribution (rough estimates)
		probs = []float64{0.15, 0.05, 0.01, 0.03, 0.76}
	}

	return utils.WeightedSample(outcomes, probs)
}
