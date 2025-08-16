package sim

import (
	"fmt"

	"github.com/logananthony/go-baseball/pkg/models"
	"github.com/logananthony/go-baseball/pkg/utils"
)

func SimulateBatterHitType(player map[string]models.BatterHitType, stand, pThrows, pitchType string, plateX, plateZ, velocity float64, teamAbbr string, pf []models.ParkFactors) string {
	zoneNum := utils.GetPitchZone(plateX, plateZ)
	veloBucket := utils.GetVelocityBucket(velocity)
	// Alias map for in-memory lookup, should match fetcher.FetchParkFactors
	abbrMap := map[string]string{
		"AZ":  "ARI",
		"ATH": "OAK",
		// "NL":  "NL", // All star game
		// etc
	}
	if val, ok := abbrMap[teamAbbr]; ok {
		teamAbbr = val
	}

	key := fmt.Sprintf("%s|%s|%d|%s|%s", stand, pThrows, zoneNum, pitchType, veloBucket)
	data, ok := player[key]
	if !ok {
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
		data.Single,
		data.Double,
		data.Triple,
		data.HomeRun,
		data.Out,
	}

	outcomes := []string{"single", "double", "triple", "home_run", "out"}

	// ─── Park Factor Adjustments ──────────────────────────────────────────
	if found {
		before := make([]float64, len(probs))
		copy(before, probs)

		probs[0] *= float64(park.OneB) / 100.0
		probs[1] *= float64(park.TwoB) / 100.0
		probs[2] *= float64(park.ThreeB) / 100.0
		probs[3] *= float64(park.HR) / 100.0
		probs[4] *= 2.0 - float64(park.H)/100.0
	} else {
		fmt.Printf("⚠️  No park factor found for team: %s\n", teamAbbr)
	}

	return utils.WeightedSample(outcomes, probs)
}
