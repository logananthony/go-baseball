package sim

import (
	"database/sql"
	"log"

	"github.com/logananthony/go-baseball/pkg/models"
)

// AggregateSprayDistributions averages the mean/std across rows, assumes fixed skew
func AggregateSprayDistributions(dists []models.SprayDistribution) models.SprayDistribution {
	if len(dists) == 0 {
		log.Fatal("No distributions provided")
	}

	var totalMean, totalStd float64
	var meanCount, stdCount int

	// Use fallback skew if first row has invalid skew
	skew := sql.NullFloat64{Float64: -0.8, Valid: true}
	if dists[0].Skew.Valid {
		skew = dists[0].Skew
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
	}
	if avgStd.Valid {
		avgStd.Float64 = totalStd / float64(stdCount)
	}

	return models.SprayDistribution{
		GameYear:          dists[0].GameYear,
		Batter:            dists[0].Batter,
		Stand:             dists[0].Stand,
		PThrows:           dists[0].PThrows,
		Outcome:           dists[0].Outcome,
		Zone:              dists[0].Zone,
		EVBucket:          dists[0].EVBucket,
		LaunchAngleBucket: dists[0].LaunchAngleBucket,
		Skew:              skew,
		Mean:              avgMean,
		Std:               avgStd,
		N:                 -1,
		Level:             "aggregated",
	}
}

func SampleFromAggregatedSprayDistribution(d models.SprayDistribution) float64 {
	if !d.Mean.Valid || !d.Std.Valid || !d.Skew.Valid {
		log.Printf("Skipping spray sample: NULL in mean/std/skew: %+v", d)
		return 0.0
	}
	return SampleFromSkewNormal(d.Mean.Float64, d.Std.Float64, d.Skew.Float64)
}
