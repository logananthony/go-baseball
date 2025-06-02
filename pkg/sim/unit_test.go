package sim

import (
	"testing"

	"github.com/logananthony/go-baseball/pkg/models"
)

func TestCachingLogic(t *testing.T) {
	// Mock input data
	simData := models.SimData{
		BatterSwing: []models.BatterSwingPercentage{
			{Batter: 1, SwingPercentage: 0.5},
			{Batter: 2, SwingPercentage: 0.6},
			{Batter: 665742, SwingPercentage: 0.7}, // BatterId 665742
		},
		BatterContact: []models.BatterContactPercentage{
			{Batter: 1, PctBallInPlay: 0.7},
			{Batter: 2, PctBallInPlay: 0.8},
			{Batter: 665742, PctBallInPlay: 0.9}, // BatterId 665742
		},
		BatterHitType: []models.BatterHitType{
			{Batter: 1, Single: 0.5},
			{Batter: 2, Double: 0.6},
			{Batter: 665742, HomeRun: 0.7}, // BatterId 665742
		},
	}

	// Populate swingPctMap
	swingPctMap := make(map[int]models.BatterSwingPercentage)
	for _, player := range simData.BatterSwing {
		swingPctMap[player.Batter] = player
	}

	// Validate swingPctMap
	t.Logf("SwingPctMap: %+v", swingPctMap)
	if swingPctMap[1].SwingPercentage != 0.5 {
		t.Errorf("Expected SwingPercentage 0.5 for Batter 1, got %f", swingPctMap[1].SwingPercentage)
	}
	if swingPctMap[2].SwingPercentage != 0.6 {
		t.Errorf("Expected SwingPercentage 0.6 for Batter 2, got %f", swingPctMap[2].SwingPercentage)
	}
	if swingPctMap[665742].SwingPercentage != 0.7 {
		t.Errorf("Expected SwingPercentage 0.7 for Batter 665742, got %f", swingPctMap[665742].SwingPercentage)
	}

	// Populate contactPctMap
	contactPctMap := make(map[int]models.BatterContactPercentage)
	for _, player := range simData.BatterContact {
		contactPctMap[player.Batter] = player
	}

	// Validate contactPctMap
	t.Logf("ContactPctMap: %+v", contactPctMap)
	if contactPctMap[1].PctBallInPlay != 0.7 {
		t.Errorf("Expected PctBallInPlay 0.7 for Batter 1, got %f", contactPctMap[1].PctBallInPlay)
	}
	if contactPctMap[2].PctBallInPlay != 0.8 {
		t.Errorf("Expected PctBallInPlay 0.8 for Batter 2, got %f", contactPctMap[2].PctBallInPlay)
	}
	if contactPctMap[665742].PctBallInPlay != 0.9 {
		t.Errorf("Expected PctBallInPlay 0.9 for Batter 665742, got %f", contactPctMap[665742].PctBallInPlay)
	}

	// Populate hitProbsMap
	hitProbsMap := make(map[int]models.BatterHitType)
	for _, player := range simData.BatterHitType {
		hitProbsMap[player.Batter] = player
	}

	// Validate hitProbsMap
	t.Logf("HitProbsMap: %+v", hitProbsMap)
	if hitProbsMap[1].Single != 0.5 {
		t.Errorf("Expected Single 0.5 for Batter 1, got %f", hitProbsMap[1].Single)
	}
	if hitProbsMap[2].Double != 0.6 {
		t.Errorf("Expected Double 0.6 for Batter 2, got %f", hitProbsMap[2].Double)
	}
	if hitProbsMap[665742].HomeRun != 0.7 {
		t.Errorf("Expected HomeRun 0.7 for Batter 665742, got %f", hitProbsMap[665742].HomeRun)
	}
}
