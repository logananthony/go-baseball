package fetcher

import (
	"os"
	"path/filepath"

	"github.com/gocarina/gocsv"
	"github.com/logananthony/go-baseball/pkg/models"
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
