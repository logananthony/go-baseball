package utils

import (
	"github.com/logananthony/go-baseball/pkg/models"
)


func GetPullProbability(probs []models.PitchingSubstitutionProb, inning, gameRuns, inningRuns int) *float64 {
	for _, prob := range probs {
		if prob.Inning == inning &&
			prob.RunsScoredGame == gameRuns &&
			prob.RunsScoredInning == inningRuns {
			
			return &prob.PullProbability
		}
	}
	return nil // no match found
}

