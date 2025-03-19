package models

type PitcherCovarianceMean struct {
    GameYear                int64
    Pitcher                 int64
    PlayerName              string
    PitchType               string
    Stand                   string
    CountState              string

    MeanPlateX              float64
    MeanPlateZ              float64
    MeanVelo                float64
    MeanPfxX                float64
    MeanPfxZ                float64

    CovPlateXPlateX         float64
    CovPlateZPlateZ         float64
    CovVeloVelo             float64
    CovPfxXPfxX             float64
    CovPfxZPfxZ             float64

    CovPlateXPlateZ         float64
    CovPlateXVelo           float64
    CovPlateZVelo           float64
    CovPfxXVelo             float64
    CovPfxZVelo             float64

    CovPfxXPlateZ           float64
    CovPfxZPlateZ           float64
    CovPfxXPlateX           float64
    CovPfxZPlateX           float64
    CovPfxXPfxZ             float64

    Count                   int64
}
