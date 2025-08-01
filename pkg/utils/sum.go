package utils

func Sum(arr []int) int {
	total := 0
	for _, x := range arr {
		total += x
	}
	return total
}
