package utils

import (
	"math"
	"sort"

	"gonum.org/v1/gonum/stat"
)

func Stddev(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	m := Mean(vals)
	var sum float64
	for _, v := range vals {
		diff := v - m
		sum += diff * diff
	}
	n := float64(len(vals))
	if n == 1 {
		return 0
	}
	// Use (n-1) for sample standard deviation
	return math.Sqrt(sum / (n - 1))
}

func Mean(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}

// QuantileWidth returns the difference between the upper and lower quantiles.
// For example, pLow=0.10 and pHigh=0.90 gives the 80% quantile width.
func QuantileWidth(data []float64, pLow, pHigh float64) float64 {
	if len(data) == 0 || pLow >= pHigh {
		return 0
	}

	// Defensive copy and sort
	sorted := make([]float64, len(data))
	copy(sorted, data)
	sort.Float64s(sorted)

	// Compute quantiles
	qLow := stat.Quantile(pLow, stat.Empirical, sorted, nil)
	qHigh := stat.Quantile(pHigh, stat.Empirical, sorted, nil)

	return qHigh - qLow
}
func BinomialCI(p float64, n int) (lower float64, upper float64) {
	if n == 0 {
		return 0, 0
	}
	z := 1.96 // 95% CI
	nF := float64(n)
	se := math.Sqrt(p * (1 - p) / nF)
	lower = p - z*se
	upper = p + z*se
	if lower < 0 {
		lower = 0
	}
	if upper > 1 {
		upper = 1
	}
	return
}

// Returns mean, lower 95%, upper 95%
func MeanCI(data []float64) (mean, lower, upper float64) {
	n := float64(len(data))
	if n == 0 {
		return 0, 0, 0
	}
	// Mean
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	mean = sum / n

	// Stddev
	sd := 0.0
	for _, v := range data {
		sd += (v - mean) * (v - mean)
	}
	std := math.Sqrt(sd / n)

	// Standard error
	se := std / math.Sqrt(n)

	// 95% CI
	z := 1.96
	lower = mean - z*se
	upper = mean + z*se
	return
}

func ProportionCI(k int, n int) (p, lower, upper float64) {
	if n == 0 {
		return 0, 0, 0
	}
	p = float64(k) / float64(n)
	se := math.Sqrt(p * (1 - p) / float64(n))
	z := 1.96 // for 95% CI
	lower = p - z*se
	upper = p + z*se
	if lower < 0 {
		lower = 0
	}
	if upper > 1 {
		upper = 1
	}
	return
}

// Usage:
// k := 390  // number of sims over 8.5
// n := 1000 // total sims
// p, lower, upper := ProportionCI(k, n)
