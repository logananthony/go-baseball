package utils

import "math"

func absInt(x int) int {
	return absDiffInt(x, 0)
}

func absDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}

func DeficitBooster(teamDeficit int, subProb float64) float64 {
	if teamDeficit < 0 && teamDeficit > -10 {
		// fmt.Println("abs teamDeficit: ", absInt(teamDeficit))
		// fmt.Println("subProb: ", subProb)
		// fmt.Println("boosted prob: ", float64(absInt(teamDeficit))*subProb)
		return math.Min(float64(absInt(teamDeficit)+1)*subProb, 1.0)
	}
	// if teamDeficit >= 0 && teamDeficit < 10 {
	// 	// fmt.Println("abs teamDeficit: ", absInt(teamDeficit))
	// 	// fmt.Println("subProb: ", subProb)
	// 	// fmt.Println("boosted prob: ", float64(absInt(teamDeficit))*subProb)
	// 	// 2
	// 	// 0.15

	// 	return (float64(absInt(teamDeficit)) * 0.01) * subProb
	// }
	return subProb
}
