package utils

import (
	//"reflect"
	"testing"
)

func ComponentTest(t *testing.T) {

  result := ReadCsvFile("~/data/batter_swing_percentage_league.csv")

  t.Logf(result)

}

