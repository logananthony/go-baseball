package utils

import (
	"math/rand"
)

var bullpenMatrix = map[int]struct {
	IDs   []int
	Probs []float64
}{
	1: {[]int{1, 2, 3}, []float64{1.0 / 3, 1.0 / 3, 1.0 / 3}},
	2: {[]int{1, 2, 3}, []float64{1.0 / 3, 1.0 / 3, 1.0 / 3}},
	3: {[]int{1, 2, 3}, []float64{1.0 / 3, 1.0 / 3, 1.0 / 3}},
	4: {[]int{3, 4, 5, 6}, []float64{0.4, 0.2, 0.2, 0.2}},
	5: {[]int{3, 4, 5, 6}, []float64{0.1, 0.4, 0.3, 0.2}},
	6: {[]int{4, 5, 6}, []float64{1.0 / 3, 1.0 / 3, 1.0 / 3}},
	7: {[]int{4, 5, 6}, []float64{1.0 / 6, 2.0 / 6, 3.0 / 6}},
	8: {[]int{5, 6, 7}, []float64{1.0 / 6, 2.0 / 6, 3.0 / 6}},
	9: {[]int{7, 8}, []float64{0.2, 0.8}},
}

func SelectBullpenPitcherLineup(pitcherLineup [][]int, inning int, used map[int]bool) []int {
	entry, ok := bullpenMatrix[inning]
	if !ok {
		return nil
	}

	var eligible [][]int
	var weights []float64

	for i, idx := range entry.IDs {
		if idx >= len(pitcherLineup) {
			continue
		}
		pid := pitcherLineup[idx][0]
		if !used[pid] {
			eligible = append(eligible, pitcherLineup[idx])
			weights = append(weights, entry.Probs[i])
		}
	}

	if len(eligible) == 0 {
		return nil
	}

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
			return eligible[i]
		}
	}

	return eligible[len(eligible)-1]
}
