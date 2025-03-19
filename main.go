package main

import (
    "github.com/logananthony/go-baseball/pkg/fetcher"
    "github.com/logananthony/go-baseball/pkg/config"
    "github.com/logananthony/go-baseball/pkg/sim"
    "fmt"
    //"github.com/logananthony/go-baseball/pkg/models"
)

func main() {
    db := config.ConnectDB()
    defer db.Close()

    
    gameYear := 2024
    batterId := 596019 
    pitcherId := 607192
    batterStands := fetcher.FetchBatterInfo(db, batterId, gameYear)
    pitcherThrows := fetcher.FetchPitcherInfo(db, pitcherId, gameYear)

    if batterStands == "B" && pitcherThrows == "R" {

      batterStands = "L"

      } else if batterStands == "B" && pitcherThrows == "L" {

      batterStands = "R"

      }

    

    pitcherFreqs := fetcher.FetchPitcherFrequencies(db, pitcherId, batterStands)
    result := sim.FilterPitcherCountPitchFreq(pitcherFreqs, 0, 0)
    //nameSlice := []string{PitcherCountPitchFreq.result.PITCH_TYPE}

    fmt.Printf("%+q", result)

  

    //for _, r := range result {
      //fmt.Printf("%+v\n", r)
    //}

  }

