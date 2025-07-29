package utils

// PitchCountToBin converts a pitch count integer to the pitch_bin string.
func PitchCountToBin(pitchCount int) string {
	switch {
	case pitchCount < 10:
		return "0-9"
	case pitchCount < 20:
		return "10-19"
	case pitchCount < 30:
		return "20-29"
	case pitchCount < 40:
		return "30-39"
	case pitchCount < 50:
		return "40-49"
	case pitchCount < 60:
		return "50-59"
	case pitchCount < 70:
		return "60-69"
	case pitchCount < 80:
		return "70-79"
	case pitchCount < 90:
		return "80-89"
	case pitchCount < 100:
		return "90-99"
	default:
		return "100+"
	}
}
