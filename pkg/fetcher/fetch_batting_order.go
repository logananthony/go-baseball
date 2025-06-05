package fetcher

import (
	"fmt"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/logananthony/go-baseball/pkg/models"
)

func FetchBattingOrder(teamAbbr string, pitcherThrows string) (models.BattingOrder, error) {
	path := "pkg/fetcher/data/batting_order.csv"
	csvFile, err := os.Open(path)
	if err != nil {
		return models.BattingOrder{}, err
	}
	defer csvFile.Close()

	var records []models.BattingOrder
	if err := gocsv.UnmarshalFile(csvFile, &records); err != nil {
		panic(err)
	}

	for _, record := range records {
		if strings.EqualFold(record.TeamAbbreviation, teamAbbr) && strings.EqualFold(record.PitcherThrows, pitcherThrows) {
			return record, nil // When match is found
		}
	}

	return models.BattingOrder{}, fmt.Errorf("no batting order found for team %s and pitcherThrows %s", teamAbbr, pitcherThrows)

}
