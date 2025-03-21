package sim

import (
	"github.com/logananthony/go-baseball/pkg/models"
  "github.com/logananthony/go-baseball/pkg/utils"
)     


// FetchPitcherFrequencies queries and prints pitch data for a pitcher ID
func SimulateBatterHitType(in []models.BatterHitType, stand, pThrows, pitchType string, plateX, plateZ, velocity float64) string {
    zone_num := utils.GetPitchZone(plateX, plateZ)
    velo_bucket := utils.GetVelocityBucket(velocity)

    hit_prob := []map[string]float64{}
    hit_weights := []int{}

    for _, each := range in {
        if each.Stand == stand || each.PThrows == pThrows || each.PitchType == pitchType || each.Zone == zone_num || each.VelocityBucket == velo_bucket {
            hit_prob = append(hit_prob, map[string]float64{
                "single":   each.Single,
                "double":   each.Double,
                "triple":   each.Triple,
                "home_run": each.HomeRun,
                "out":      each.Out,
            })
            hit_weights = append(hit_weights, each.N)
        }
    }

    // Aggregate into final smoothed probabilities
    smoothed := map[string]float64{"single": 0, "double": 0, "triple": 0, "home_run": 0, "out": 0}
    total_weight := 0

    for i, probs := range hit_prob {
        n := hit_weights[i]
        total_weight += n
        for k, v := range probs {
            smoothed[k] += v * float64(n)
        }
    }

    // Normalize
    for k := range smoothed {
        smoothed[k] /= float64(total_weight)
    }

    outcomes := []string{"single", "double", "triple", "home_run", "out"}
    probs := []float64{
        smoothed["single"],
        smoothed["double"],
        smoothed["triple"],
        smoothed["home_run"],
        smoothed["out"],
    }

    // Sample from smoothed
    return utils.WeightedSample(outcomes, probs) 
}

