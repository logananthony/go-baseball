package sim

import (
	"github.com/logananthony/go-baseball/pkg/models"
	//"github.com/logananthony/go-baseball/pkg/utils"
	"time"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distmv"
	//"gonum.org/v1/gonum/mat"
)

// FetchPitcherFrequencies queries and prints pitch data for a pitcher ID
func SimulatePitchLocationVelo(player []models.PitcherCovarianceMean, league []models.PitcherCovarianceMeanLeague, pitch_type string, stand string, balls, strikes int) []float64  {

    var count_state string
    if balls == strikes {
       count_state = "even"
    } else if balls > strikes {
       count_state = "behind"
      } else {
       count_state = "ahead"
    }

    src := rand.NewSource(uint64(time.Now().UnixNano()))


   player_mean_mat := []float64{}
   player_cov_mat := mat.NewSymDense(3, []float64 {
        0, 0, 0,
        0, 0, 0,
        0, 0, 0,
        })

        for _, each := range player {
          if each.CountState == count_state && each.PitchType == pitch_type && each.Stand == stand {

            if each.Count >= 30 {
              
              player_mean_mat = append(player_mean_mat, each.MeanPlateX, each.MeanPlateZ, each.MeanVelo)

              player_cov_mat = mat.NewSymDense(3, []float64 {
                  each.CovPlateXPlateX, each.CovPlateXPlateZ, each.CovPlateXVelo,
                  each.CovPlateXPlateZ, each.CovPlateZPlateZ, each.CovPlateZVelo,
                  each.CovPlateXVelo, each.CovPlateZVelo, each.CovVeloVelo,
              })

              break 

            }
          }
      } 

      if player_mean_mat == nil {
          sample := distmv.NormalRandCov(nil, player_mean_mat, player_cov_mat, src)
          return sample
      } 

   league_mean_mat := []float64{}
   league_cov_mat := mat.NewSymDense(3, []float64 {
        0, 0, 0,
        0, 0, 0,
        0, 0, 0,
        })

        for _, each := range league {
          if each.CountState == count_state && each.PitchType == pitch_type && each.Stand == stand {
              
              league_mean_mat = append(league_mean_mat, each.MeanPlateX, each.MeanPlateZ, each.MeanVelo)

              league_cov_mat = mat.NewSymDense(3, []float64 {
                  each.CovPlateXPlateX, each.CovPlateXPlateZ, each.CovPlateXVelo,
                  each.CovPlateXPlateZ, each.CovPlateZPlateZ, each.CovPlateZVelo,
                  each.CovPlateXVelo, each.CovPlateZVelo, each.CovVeloVelo,
              })
              sample := distmv.NormalRandCov(nil, league_mean_mat, league_cov_mat, src)
              return sample

              
          }
      } 


    return []float64{0, 0, 0}


}



