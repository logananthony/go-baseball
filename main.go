package main

import (
    //"github.com/logananthony/go-baseball/pkg/fetcher"
    "github.com/logananthony/go-baseball/pkg/models"
    "github.com/logananthony/go-baseball/pkg/sim"
    "github.com/davecgh/go-spew/spew"
    //"fmt"
    //"encoding/json"
)

func main() {

  sim_results := sim.SimulateAtBat([]models.AtBatData{
    {
        GameYear: 2024,
        PitcherId: 607192,
        BatterId: 592450,
        Strikes: 0,
        Balls: 0,
    },
  })

  spew.Dump(sim_results)


  }

