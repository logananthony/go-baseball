package sim

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/logananthony/go-baseball/pkg/models"
)

// AggregateEVDistributions averages the mean/std across rows, assumes fixed skew
func AggregateEVDistributions(dists []models.EVDistribution) models.EVDistribution {
	if len(dists) == 0 {
		log.Fatal("No distributions provided")
	}

	var totalMean, totalStd float64
	skew := dists[0].Skew // Assume all skew values are the same
	count := len(dists)

	for _, d := range dists {
		totalMean += d.Mean
		totalStd += d.Std
	}

	return models.EVDistribution{
		GameYear:       dists[0].GameYear,
		Batter:         dists[0].Batter,
		Stand:          dists[0].Stand,
		PThrows:        dists[0].PThrows,
		Outcome:        dists[0].Outcome,
		PitchType:      dists[0].PitchType,
		Zone:           dists[0].Zone,
		VelocityBucket: dists[0].VelocityBucket,
		Skew:           skew,
		Mean:           totalMean / float64(count),
		Std:            totalStd / float64(count),
		N:              -1,
		Level:          "aggregated",
	}
}

// SampleFromSkewNormal returns a sample from a skew-normal distribution using Azzalini's method
func SampleFromSkewNormal(mean, std, skew float64) float64 {
	rand.Seed(time.Now().UnixNano())

	// Generate two standard normals
	u0 := rand.NormFloat64()
	v := rand.NormFloat64()

	// Compute delta from skew
	delta := skew / math.Sqrt(1+skew*skew)

	// Create Z ~ SkewNormal(0, 1, skew)
	z := delta*u0 + math.Sqrt(1-delta*delta)*v
	sample := mean + std*z
	return sample
}

// SampleFromAggregatedDistribution draws using the aggregated parameters
func SampleFromAggregatedDistribution(d models.EVDistribution) float64 {
	return SampleFromSkewNormal(d.Mean, d.Std, d.Skew)
}

