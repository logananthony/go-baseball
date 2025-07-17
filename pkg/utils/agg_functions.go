package utils

import "math"

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
