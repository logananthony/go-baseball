package sim

import (
	"log"
	"github.com/logananthony/go-baseball/pkg/models"
)

// AggregateEVDistributions averages the mean/std across rows, assumes fixed skew
func AggregateLADistributions(dists []models.LADistribution) models.LADistribution {
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

	return models.LADistribution{
		GameYear:       dists[0].GameYear,
		Batter:         dists[0].Batter,
		Stand:          dists[0].Stand,
		PThrows:        dists[0].PThrows,
		Outcome:        dists[0].Outcome,
		Zone:           dists[0].Zone,
		EVBucket:       dists[0].EVBucket,
		Skew:           skew,
		Mean:           totalMean / float64(count),
		Std:            totalStd / float64(count),
		N:              -1,
		Level:          "aggregated",
	}
}

func SampleFromAggregatedLADistribution(d models.LADistribution) float64 {
	return SampleFromSkewNormal(d.Mean, d.Std, d.Skew)
}


