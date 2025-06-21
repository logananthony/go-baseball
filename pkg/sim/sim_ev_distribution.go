package sim

import (
	"database/sql"
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
	var meanCount, stdCount int

	skew := sql.NullFloat64{Float64: -0.8, Valid: true} // default fallback
	for _, d := range dists {
		if d.Skew.Valid {
			skew = d.Skew
			break
		}
	}

	for _, d := range dists {
		if d.Mean.Valid {
			totalMean += d.Mean.Float64
			meanCount++
		}
		if d.Std.Valid {
			totalStd += d.Std.Float64
			stdCount++
		}
	}

	avgMean := sql.NullFloat64{Valid: meanCount > 0}
	avgStd := sql.NullFloat64{Valid: stdCount > 0}

	if avgMean.Valid {
		avgMean.Float64 = totalMean / float64(meanCount)
	} else {
		log.Println("No valid mean found; returning 0")
	}
	if avgStd.Valid {
		avgStd.Float64 = totalStd / float64(stdCount)
	} else {
		log.Println("No valid std found; returning 0")
	}

	// Use first row as template for identifiers (or fallback to nulls)
	template := models.EVDistribution{}
	for _, d := range dists {
		if d.Batter != 0 {
			template = d
			break
		}
	}

	return models.EVDistribution{
		GameYear:       template.GameYear,
		Batter:         template.Batter,
		Stand:          template.Stand,
		PThrows:        template.PThrows,
		Outcome:        template.Outcome,
		PitchType:      template.PitchType,
		Zone:           template.Zone,
		VelocityBucket: template.VelocityBucket,
		Skew:           skew,
		Mean:           avgMean,
		Std:            avgStd,
		N:              sql.NullInt32{Int32: int32(len(dists)), Valid: true},
		Level:          "aggregated",
	}
}

// SampleFromSkewNormal returns a sample from a skew-normal distribution using Azzalini's method
func SampleFromSkewNormal(mean, std, skew float64) float64 {
	rand.Seed(time.Now().UnixNano())

	u0 := rand.NormFloat64()
	v := rand.NormFloat64()

	delta := skew / math.Sqrt(1+skew*skew)
	z := delta*u0 + math.Sqrt(1-delta*delta)*v

	return mean + std*z
}

// SampleFromAggregatedDistribution draws using the aggregated parameters
func SampleFromAggregatedDistribution(d models.EVDistribution) float64 {
	if !d.Mean.Valid || !d.Std.Valid || !d.Skew.Valid {
		log.Printf("Skipping sample: missing Mean/Std/Skew in %+v", d)
		return 0.0
	}

	return SampleFromSkewNormal(d.Mean.Float64, d.Std.Float64, d.Skew.Float64)
}
