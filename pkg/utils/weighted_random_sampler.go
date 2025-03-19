package utils

import (
    "math/rand"
    //"time"
)

func WeightedSample(choices []string, weights []float64) string {
    //rand.Seed(time.Now().UnixNano())

    // Compute cumulative sum of weights
    var cumulative []float64
    total := 0.0
    for _, w := range weights {
        total += w
        cumulative = append(cumulative, total)
    }

    // Generate a random number between 0 and total weight
    r := rand.Float64() * total

    // Find the corresponding choice
    for i, c := range cumulative {
        if r < c {
            return choices[i]
        }
    }

    // Fallback (shouldn't happen)
    return choices[len(choices)-1]
}

