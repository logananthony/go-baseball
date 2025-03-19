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
func SimulatePitchLocation(in []models.PitcherCovarianceMean, pitch_type string, stand string, balls, strikes int) []float64  {

    var count_state string
    if balls == strikes {
       count_state = "even"
    } else if balls > strikes {
       count_state = "behind"
      } else {
       count_state = "ahead"
    }

   mean_mat := []float64{}
   //var cov_mat [3][3]float64
   cov_mat := mat.NewSymDense(3, []float64 {
        0, 0, 0,
        0, 0, 0,
        0, 0, 0,
        })

        for _, each := range in {
          if each.CountState == count_state && each.PitchType == pitch_type && each.Stand == stand {
              
              mean_mat = append(mean_mat, each.MeanPlateX, each.MeanPlateZ, each.MeanVelo)

              
              cov_mat = mat.NewSymDense(3, []float64 {
                  each.CovPlateXPlateX, each.CovPlateXPlateZ, each.CovPlateXVelo,
                  each.CovPlateXPlateZ, each.CovPlateZPlateZ, each.CovPlateZVelo,
                  each.CovPlateXVelo, each.CovPlateZVelo, each.CovVeloVelo,
              })
              break 
          }
      } 

    src := rand.NewSource(uint64(time.Now().UnixNano()))


    


    sample := distmv.NormalRandCov(nil, mean_mat, cov_mat, src)
    



  return sample

}



