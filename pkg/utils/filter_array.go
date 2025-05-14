package utils

func FilterSliceSlices(sli [][]int, search int) [][]int {
	for i, v := range sli {
		if v[0] == search {
			return append(sli[:i], sli[i+1:]...)
		}
	}
	return sli

}
