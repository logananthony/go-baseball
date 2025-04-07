package utils

import (
    "math/rand"
    //"time"
)

func WeightedSample(choices []string, weights []float64) string {
    if len(choices) == 0 || len(weights) == 0 || len(choices) != len(weights) {
        return "" // or log, panic with message, or return a default value
    }

    var cumulative []float64
    total := 0.0
    for _, w := range weights {
        total += w
        cumulative = append(cumulative, total)
    }

    if total == 0.0 {
        return choices[0] // or handle gracefully
    }

    r := rand.Float64() * total

    for i, c := range cumulative {
        if r < c {
            return choices[i]
        }
    }

    return choices[len(choices)-1] // fallback
}

