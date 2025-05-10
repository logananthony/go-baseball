package utils

import (
	"math/rand"
	"golang.org/x/exp/constraints"
)

func WeightedSample[T any, W constraints.Integer | constraints.Float](choices []T, weights []W) T {
	var zero T // default zero value for T

	if len(choices) == 0 || len(weights) == 0 || len(choices) != len(weights) {
		return zero
	}

	var cumulative []float64
	total := 0.0
	for _, w := range weights {
		total += float64(w) // cast weight to float64
		cumulative = append(cumulative, total)
	}

	if total == 0.0 {
		return choices[0]
	}

	r := rand.Float64() * total

	for i, c := range cumulative {
		if r < c {
			return choices[i]
		}
	}

	return choices[len(choices)-1]
}

