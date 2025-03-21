package utils

func GetVelocityBucket(velocity float64) string {
	switch {
	case velocity < 70:
		return "<70"
	case velocity >= 70 && velocity < 75:
		return "70-75"
	case velocity >= 75 && velocity < 80:
		return "75-80"
	case velocity >= 80 && velocity < 85:
		return "80-85"
	case velocity >= 85 && velocity < 90:
		return "85-90"
	case velocity >= 90 && velocity < 95:
		return "90-95"
	case velocity >= 95 && velocity < 100:
		return "95-100"
	default:
		return "100+"
	}
}

