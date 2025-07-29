package utils

import (
	"math/rand"

	"github.com/logananthony/go-baseball/pkg/models"
)

// SelectBullpenPitcherLineup chooses a bullpen pitcher based on situation and not-yet-used pitchers.
// Returns a pitcher [playerID, season] pair or nil if no eligible pitcher found.
func SelectBullpenPitcherLineup(
	pitcherLineup [][]int,
	roleProbs []models.BullpenRoleProb,
	inning int,
	runsScoredGame int,
	runsScoredInning int,
	used map[int]bool,
) []int {
	// Filter for relevant role probabilities for this situation
	var eligibleRoles []int
	var weights []float64

	for _, rp := range roleProbs {
		if rp.Inning == inning &&
			rp.RunsScoredGame == runsScoredGame &&
			rp.RunsScoredInning == runsScoredInning {
			if rp.Role >= len(pitcherLineup) {
				continue
			}
			pid := pitcherLineup[rp.Role][0]
			if !used[pid] {
				eligibleRoles = append(eligibleRoles, rp.Role)
				weights = append(weights, rp.PullProbability)
			}
		}
	}
	if len(eligibleRoles) == 0 {
		return nil
	}
	// Normalize weights and select randomly
	total := 0.0
	for _, w := range weights {
		total += w
	}
	for i := range weights {
		weights[i] /= total
	}
	r := rand.Float64()
	cumulative := 0.0
	for i, prob := range weights {
		cumulative += prob
		if r <= cumulative {
			return pitcherLineup[eligibleRoles[i]]
		}
	}
	return pitcherLineup[eligibleRoles[len(eligibleRoles)-1]]
}
