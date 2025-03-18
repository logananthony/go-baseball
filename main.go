package main

import (
    "github.com/logananthony/go-baseball/pkg/fetcher"
    "github.com/logananthony/go-baseball/pkg/config"
)

func main() {
    db := config.ConnectDB() // Call it from config, not fetcher
    defer db.Close()

    //var batter_info =
    gameYear := 2024
    batterId := 596019 
    pitcherId := 607192
    batterStands := fetcher.FetchBatterInfo(db, batterId, gameYear)
    pitcherThrows := fetcher.FetchPitcherInfo(db, pitcherId, gameYear)

    //var batterStands string

    if batterStands == "B" && pitcherThrows == "R" {

      batterStands = "L"

      } else if batterStands == "B" && pitcherThrows == "L" {

      batterStands = "R"

      } 

    fetcher.FetchPitcherFrequencies(db, pitcherId, batterStands)
    //println(pitcherInfo)
}

