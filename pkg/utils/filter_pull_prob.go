package utils

import (
	"github.com/logananthony/go-baseball/pkg/models"
)

func GetPitcherPulled(probs []models.BullpenRoleProb, inning, gameRuns, inningRuns int) *float64 {
	for i := range probs {
		if probs[i].Inning == inning &&
			probs[i].RunsScoredGame == gameRuns &&
			probs[i].RunsScoredInning == inningRuns {

			return &probs[i].PullProbability
		}
	}
	return nil
}
