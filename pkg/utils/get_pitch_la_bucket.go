package utils

import "fmt"

// GetLaunchAngleBucket returns the 10-degree bucket for a given launch angle.
// Example: 17.3 -> "10-20"
func GetLaunchAngleBucket(launchAngle float64) string {
	lower := int(launchAngle) / 10 * 10
	upper := lower + 10
	return fmt.Sprintf("%d-%d", lower, upper)
}

