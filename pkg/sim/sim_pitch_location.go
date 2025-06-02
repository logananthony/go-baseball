package sim

// import (
// 	"fmt"
// 	"math/rand"
// 	"time"

// 	"github.com/logananthony/go-baseball/pkg/models"
// 	"gonum.org/v1/gonum/mat"
// 	"gonum.org/v1/gonum/stat/distmv"
// )

//"github.com/logananthony/go-baseball/pkg/utils"

//"fmt"
//"gonum.org/v1/gonum/mat"

// FetchPitcherFrequencies queries and prints pitch data for a pitcher ID
// func SimulatePitchLocationVelo(player []models.PitcherCovarianceMean, league []models.PitcherCovarianceMeanLeague, pitch_type string, stand string, balls, strikes int) []float64 {

// 	var count_state string
// 	if balls == strikes {
// 		count_state = "even"
// 	} else if balls > strikes {
// 		count_state = "behind"
// 	} else {
// 		count_state = "ahead"
// 	}

// 	src := rand.NewSource(uint64(time.Now().UnixNano()))

// 	player_mean_mat := []float64{}
// 	var player_cov_mat *mat.SymDense

// 	for _, each := range player {
// 		if each.CountState == count_state && each.PitchType == pitch_type && each.Stand == stand {

// 			if each.Count >= 1 {

// 				player_mean_mat = append(player_mean_mat, each.MeanPlateX, each.MeanPlateZ, each.MeanVelo)

// 				player_cov_mat = mat.NewSymDense(3, []float64{
// 					each.CovPlateXPlateX, each.CovPlateXPlateZ, each.CovPlateXVelo,
// 					each.CovPlateXPlateZ, each.CovPlateZPlateZ, each.CovPlateZVelo,
// 					each.CovPlateXVelo, each.CovPlateZVelo, each.CovVeloVelo,
// 				})

// 				break

// 			}
// 		}
// 	}

// 	if len(player_mean_mat) > 0 && player_cov_mat != nil {

// 		sample := distmv.NormalRandCov(nil, player_mean_mat, player_cov_mat, src)
// 		// fmt.Println(sample)
// 		return sample
// 	}

// 	league_mean_mat := []float64{}
// 	var league_cov_mat *mat.SymDense

// 	for _, each := range league {
// 		if each.CountState == count_state && each.PitchType == pitch_type && each.Stand == stand {

// 			league_mean_mat = append(league_mean_mat, each.MeanPlateX, each.MeanPlateZ, each.MeanVelo)

// 			league_cov_mat = mat.NewSymDense(3, []float64{
// 				each.CovPlateXPlateX, each.CovPlateXPlateZ, each.CovPlateXVelo,
// 				each.CovPlateXPlateZ, each.CovPlateZPlateZ, each.CovPlateZVelo,
// 				each.CovPlateXVelo, each.CovPlateZVelo, each.CovVeloVelo,
// 			})
// 			sample := distmv.NormalRandCov(nil, league_mean_mat, league_cov_mat, src)
// 			return sample

// 		}
// 	}

// 	return []float64{0, 0, 0}

// }

import (
	"time"

	"github.com/logananthony/go-baseball/pkg/models"
	"golang.org/x/exp/rand" // Correct package
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distmv"
)

func SimulatePitchLocationVelo(player []models.PitcherCovarianceMean, league []models.PitcherCovarianceMeanLeague, pitch_type string, stand string, balls, strikes int) []float64 {
	var count_state string
	if balls == strikes {
		count_state = "even"
	} else if balls > strikes {
		count_state = "behind"
	} else {
		count_state = "ahead"
	}

	// fmt.Printf("SimulatePitchLocationVelo called with:\n")
	// fmt.Printf("  Pitch Type: %s, Stand: %s, Balls: %d, Strikes: %d, Count State: %s\n", pitch_type, stand, balls, strikes, count_state)

	src := rand.NewSource(uint64(time.Now().UnixNano())) // Use uint64 for Seed
	player_mean_mat := []float64{}
	var player_cov_mat *mat.SymDense

	// Check player data
	// fmt.Println("Checking player data...")
	for _, each := range player {
		// fmt.Printf("  Player Entry: %+v\n", each)
		if each.CountState == count_state && each.PitchType == pitch_type && each.Stand == stand {
			// fmt.Println("  Match found in player data.")
			if each.Count >= 1 {
				// fmt.Println("  Player data has sufficient count.")
				player_mean_mat = append(player_mean_mat, each.MeanPlateX, each.MeanPlateZ, each.MeanVelo)

				player_cov_mat = mat.NewSymDense(3, []float64{
					each.CovPlateXPlateX, each.CovPlateXPlateZ, each.CovPlateXVelo,
					each.CovPlateXPlateZ, each.CovPlateZPlateZ, each.CovPlateZVelo,
					each.CovPlateXVelo, each.CovPlateZVelo, each.CovVeloVelo,
				})
				break
			} else {
				// fmt.Println("  Player data does not have sufficient count.")
			}
		}
	}

	if len(player_mean_mat) > 0 && player_cov_mat != nil {
		// fmt.Println("Player data successfully used for simulation.")
		sample := distmv.NormalRandCov(nil, player_mean_mat, player_cov_mat, src)
		return sample
	} else {
		// fmt.Println("No valid player data found.")
	}

	league_mean_mat := []float64{}
	var league_cov_mat *mat.SymDense

	// Check league data
	// fmt.Println("Checking league data...")
	for _, each := range league {
		// fmt.Printf("  League Entry: %+v\n", each)
		if each.CountState == count_state && each.PitchType == pitch_type && each.Stand == stand {
			// fmt.Println("  Match found in league data.")
			league_mean_mat = append(league_mean_mat, each.MeanPlateX, each.MeanPlateZ, each.MeanVelo)

			league_cov_mat = mat.NewSymDense(3, []float64{
				each.CovPlateXPlateX, each.CovPlateXPlateZ, each.CovPlateXVelo,
				each.CovPlateXPlateZ, each.CovPlateZPlateZ, each.CovPlateZVelo,
				each.CovPlateXVelo, each.CovPlateZVelo, each.CovVeloVelo,
			})
			sample := distmv.NormalRandCov(nil, league_mean_mat, league_cov_mat, src)
			return sample
		}
	}

	if len(league_mean_mat) == 0 || league_cov_mat == nil {
		// fmt.Println("No valid league data found.")
	}

	// fmt.Println("Returning default value: [0, 0, 0]")
	return []float64{0, 0, 0}
}
