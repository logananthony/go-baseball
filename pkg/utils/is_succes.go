package utils

import(
  "math/rand"
  "math"
)

func IsSuccess(prob_success *float64) bool {
    p := *prob_success
    hit := rand.Float64() < p
    if hit {
        p = math.Max(0, p - 0.1)
    } else {
        p = math.Min(1, p + 0.1)
    }
    *prob_success = p
    return hit
}
