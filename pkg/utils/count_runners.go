package utils

// CountRunners returns the number of true values in the first 3 bases
func CountRunners(baseState []bool) int {
	count := 0
	for i := 0; i < 3 && i < len(baseState); i++ {
		if baseState[i] {
			count++
		}
	}
	return count
}
