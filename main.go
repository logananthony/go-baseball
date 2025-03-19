package main

import (
    "github.com/logananthony/go-baseball/pkg/fetcher"
    "github.com/logananthony/go-baseball/pkg/config"
    "github.com/logananthony/go-baseball/pkg/sim"
    "fmt"
    //"encoding/json"
)

func main() {

    db := config.ConnectDB()
    defer db.Close()

    gameYear := 2024
    batterId := 596019 
    pitcherId := 607192
    batterStands := fetcher.FetchBatterInfo(db, batterId, gameYear)
    pitcherThrows := fetcher.FetchPitcherInfo(db, pitcherId, gameYear)
    strikeCount := 0
    ballCount := 0

    if batterStands == "B" && pitcherThrows == "R" {
      batterStands = "L"
      } else if batterStands == "B" && pitcherThrows == "L" {
      batterStands = "R"
      }

    pitcherFreqs := fetcher.FetchPitcherFrequencies(db, pitcherId, batterStands)
    pitch_result := sim.SimulatePitchType(pitcherFreqs, ballCount, strikeCount)

    //fmt.Printf("%+q", result)

    covariance := fetcher.FetchPitcherCovarianceMean(db, int64(pitcherId), int64(gameYear))
    //s, _ := json.MarshalIndent(covariance, "", "\t")
    //fmt.Print(string(s))



    location := sim.SimulatePitchLocation(covariance, pitch_result, batterStands, ballCount, strikeCount)

    fmt.Println(location)


  }

