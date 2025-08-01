package utils

// CountRunners returns the number of occupied bases in the first 3 spots
func CountRunners(baseState []int) int {
	count := 0
	for i := 0; i < 3 && i < len(baseState); i++ {
		if baseState[i] != 0 {
			count++
		}
	}
	return count
}
