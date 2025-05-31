package models

type PlateAppearanceData struct {
	BatterGameYear  int
	BatterId        int
	PitcherGameYear int
	PitcherId       int
	Strikes         int
	Balls           int
	AwayScore       int
	HomeScore       int
	AtBatNumber     int
	Inning          int
	InningTopBot    string
	Outs            int
	On1b            bool
	On2b            bool
	On3b            bool
}
