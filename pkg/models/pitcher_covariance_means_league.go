package models

type PitcherCovarianceMeanLeague struct {
	PitchType            string  `csv:"pitch_type"`
	Stand                string  `csv:"stand"`
	CountState           string  `csv:"count_state"`
	MeanPlateX           float64 `csv:"mean_plate_x"`
	MeanPlateZ           float64 `csv:"mean_plate_z"`
	MeanVelo             float64 `csv:"mean_velo"`
	MeanPfxX             float64 `csv:"mean_pfx_x"`
	MeanPfxZ             float64 `csv:"mean_pfx_z"`
	CovPlateXPlateX      float64 `csv:"cov_plate_x_platex"`
	CovPlateZPlateZ      float64 `csv:"cov_plate_z_plate_z"`
	CovVeloVelo          float64 `csv:"cov_velo_velo"`
	CovPfxXPfxX          float64 `csv:"cov_pfx_x_pfx_x"`
	CovPfxZPfxZ          float64 `csv:"cov_pfx_z_pfx_z"`
	CovPlateXPlateZ      float64 `csv:"cov_plate_x_plate_z"`
	CovPlateXVelo        float64 `csv:"cov_plate_x_velo"`
	CovPlateZVelo        float64 `csv:"cov_plate_z_velo"`
	CovPfxXVelo          float64 `csv:"cov_pfx_x_velo"`
	CovPfxZVelo          float64 `csv:"cov_pfx_z_velo"`
	CovPfxXPlateZ        float64 `csv:"cov_pfx_x_plate_z"`
	CovPfxZPlateZ        float64 `csv:"cov_pfx_z_plate_z"`
	CovPfxXPlateX        float64 `csv:"cov_pfx_x_plate_x"`
	CovPfxZPlateX        float64 `csv:"cov_pfx_z_plate_x"`
	CovPfxXPfxZ          float64 `csv:"cov_pfx_x_pfx_z"`
	Count                int     `csv:"count"`
}

