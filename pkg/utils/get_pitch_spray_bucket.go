package utils

import "math"

func GetSprayBucket(angle float64) string {
	if math.IsNaN(angle) {
		return ""
	}
	if angle <= -20 {
		return "-45 to -20"
	} else if angle > -20 && angle <= -5 {
		return "-20 to -5"
	} else if angle > -5 && angle <= 5 {
		return "-5 to 5"
	} else if angle > 5 && angle <= 20 {
		return "5 to 20"
	} else {
		return "20 to 45"
	}
}

