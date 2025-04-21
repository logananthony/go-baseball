package utils

import "fmt"

func GetEVBucket(ev float64) string {
	if ev < 0 {
		return ""
	}
	start := int(ev) / 5 * 5
	end := start + 5
	return fmt.Sprintf("%d-%d", start, end)
}

