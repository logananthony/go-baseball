package sim

import (
	"github.com/logananthony/go-baseball/pkg/models"
	"github.com/logananthony/go-baseball/pkg/utils"
	// "fmt"
)

// FetchPitcherFrequencies queries and prints pitch data for a pitcher ID
func SimulateContactPercentage(player []models.BatterContactPercentage, league []models.BatterContactPercentageLeague, stand, pThrows, pitchType string, plateX, plateZ float64) string  {
    
    zoneNum := utils.GetPitchZone(plateX, plateZ)
    //fmt.Println(zoneNum, pitchType)
    
    pitch_types := []string{"swinging_strike", "foul", "ball_in_play"} 
    player_contact_prob := []float64{}
    league_contact_prob := []float64{}
    var contactResult string

    for _, each := range player {
      if each.Stand == stand && each.PThrows == pThrows && each.PitchType == pitchType && each.Zone == zoneNum {
        if each.TotalSwings >= 25 {
          player_contact_prob = append(player_contact_prob, each.PctSwingingStrike, each.PctFoul, each.PctBallInPlay)
          //fmt.Println(player_contact_prob)
          break
        }
      }
    }


    for _, each := range league {
      if each.Stand == stand && each.PThrows == pThrows && each.PitchType == pitchType && each.Zone == zoneNum {
          league_contact_prob = append(league_contact_prob, each.PctSwingingStrike, each.PctFoul, each.PctBallInPlay)
        //return utils.WeightedSample(pitch_types, league_contact_prob)
      }
    }

   if len(player_contact_prob) == len(pitch_types) {
    contactResult = utils.WeightedSample(pitch_types, player_contact_prob)
    } else if len(league_contact_prob) == len(pitch_types) {
        contactResult = utils.WeightedSample(pitch_types, league_contact_prob)
    } else {
        contactResult = "ball_in_play" // or "no_data", or return an error
    }

    

    return contactResult

}




