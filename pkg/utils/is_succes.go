package utils

import (
	"math"
	"math/rand"
)

// IsSuccess randomly returns true with probability *probSuccess,
// and updates *probSuccess by decreasing on success and increasing on failure.
// If probSuccess is nil, always returns false.
func IsSuccess(probSuccess *float64) bool {
	if probSuccess == nil {
		return false
	}

	p := *probSuccess
	hit := rand.Float64() < p

	if hit {
		p = math.Max(0, p-0.1)
	} else {
		p = math.Min(1, p+0.1)
	}
	*probSuccess = p

	return hit
}
