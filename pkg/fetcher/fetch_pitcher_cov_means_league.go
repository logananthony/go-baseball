package fetcher

import (
    "path/filepath"
    "os"
    "github.com/logananthony/go-baseball/pkg/models"
    "github.com/gocarina/gocsv"
)

func FetchPitcherCovarianceMeanLeague() []models.PitcherCovarianceMeanLeague {
    cwd, _ := os.Getwd()
    path := filepath.Join(cwd, "pkg", "fetcher", "data", "pitcher_covariance_mean_league.csv")
    csvFile, err := os.Open(path)
    if err != nil {
        panic(err)
    }
    defer csvFile.Close()

    var records []models.PitcherCovarianceMeanLeague
    if err := gocsv.UnmarshalFile(csvFile, &records); err != nil {
        panic(err)
    }

    return records
}

