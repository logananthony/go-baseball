package models

type SimOutcome struct {
	HomeScore     int
	AwayScore     int
	GamePk        int
	BatterEvents  []BatterEvent
	PitcherEvents []PitcherEvent
	RunEvents     []RunEvent
}
