package fetcher

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/logananthony/go-baseball/pkg/models"
)

func FetchBullpenOrder(teamAbbr string) *models.BullpenOrder {
	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, "pkg", "fetcher", "data", "bullpen_order.csv")
	csvFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	var records []models.BullpenOrder
	if err := gocsv.UnmarshalFile(csvFile, &records); err != nil {
		panic(err)
	}

	for _, record := range records {
		if strings.EqualFold(record.TeamAbbreviation, teamAbbr) {
			return &record
		}
	}

	return nil // No match found
}
