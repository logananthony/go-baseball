package fetcher

import (
    "path/filepath"
    "os"
    "github.com/logananthony/go-baseball/pkg/models"
    "github.com/gocarina/gocsv"
)

func FetchBatterContactPercentageLeague() []models.BatterContactPercentageLeague {
    cwd, _ := os.Getwd()
    path := filepath.Join(cwd, "pkg", "fetcher", "data", "batter_contact_percentages_league.csv")
    csvFile, err := os.Open(path)
    if err != nil {
        panic(err)
    }
    defer csvFile.Close()

    var records []models.BatterContactPercentageLeague
    if err := gocsv.UnmarshalFile(csvFile, &records); err != nil {
        panic(err)
    }

    return records
}

