package utils

func ConvertPitcherThrows(pitcherThrows *string) *string {
	if pitcherThrows != nil {
		switch *pitcherThrows {
		case "Left":
			left := "L"
			return &left
		case "Right":
			right := "R"
			return &right
		}
	}
	return pitcherThrows
}
