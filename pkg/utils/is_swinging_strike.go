package utils

func IsSwingingStrikeCount(is_swing, is_strike []bool) int {
	n := len(is_swing)
	if len(is_strike) < n {
		n = len(is_strike)
	}

	total := 0
	for i := 0; i < n; i++ {
		if is_swing[i] && is_strike[i] {
			total++
		}
	}
	return total
}

func CountTrueBool(booleanSlice []bool) int {
	n := len(booleanSlice)
	total := 0
	for i := 0; i < n; i++ {
		if booleanSlice[i] {
			total++
		}
	}
	return total
}
