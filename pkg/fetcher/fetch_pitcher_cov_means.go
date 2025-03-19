package fetcher

import (
	"database/sql"
	//"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/logananthony/go-baseball/pkg/models"
)



// FetchPitcherFrequencies queries and prints pitch data for a pitcher ID
func FetchPitcherCovarianceMean(db *sql.DB, pitcherID int64, gameYear int64) []models.PitcherCovarianceMean {
    query := `SELECT game_year, pitcher, player_name, pitch_type, stand, count_state, mean_plate_x, mean_plate_z, mean_velo,
                      mean_pfx_x, mean_pfx_z, cov_plate_x_platex, cov_plate_z_plate_z, cov_velo_velo, cov_pfx_x_pfx_x, 
                      cov_pfx_z_pfx_z, cov_plate_x_plate_z, cov_plate_x_velo, cov_plate_z_velo, cov_pfx_x_velo, cov_pfx_z_velo, 
                      cov_pfx_x_plate_z, cov_pfx_z_plate_z, cov_pfx_x_plate_x, cov_pfx_z_plate_x, cov_pfx_x_pfx_z, 
                      count FROM pitcher_covariance_mean WHERE pitcher = $1 AND game_year = $2`

    rows, err := db.Query(query, pitcherID, gameYear)
    if err != nil {
        log.Fatal("Query error:", err)
    }
    defer rows.Close()

    var results []models.PitcherCovarianceMean

    for rows.Next() {
        var freq models.PitcherCovarianceMean
        freq.GameYear = gameYear
        freq.Pitcher = pitcherID 
        err := rows.Scan(
            &freq.GameYear,
            &freq.Pitcher,
            &freq.PlayerName,
            &freq.PitchType,
            &freq.Stand,
            &freq.CountState,

            &freq.MeanPlateX,
            &freq.MeanPlateZ,
            &freq.MeanVelo,
            &freq.MeanPfxX,
            &freq.MeanPfxZ,

            &freq.CovPlateXPlateX,
            &freq.CovPlateZPlateZ,
            &freq.CovVeloVelo,
            &freq.CovPfxXPfxX,
            &freq.CovPfxZPfxZ,

            &freq.CovPlateXPlateZ,
            &freq.CovPlateXVelo,
            &freq.CovPlateZVelo,
            &freq.CovPfxXVelo,
            &freq.CovPfxZVelo,

            &freq.CovPfxXPlateZ,
            &freq.CovPfxZPlateZ,
            &freq.CovPfxXPlateX,
            &freq.CovPfxZPlateX,
            &freq.CovPfxXPfxZ,

            &freq.Count,
        )
        if err != nil {
            log.Fatal(err)
        }

        results = append(results, freq)
    }

    return results
}


