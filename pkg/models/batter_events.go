package models

type BatterEvent struct {
	BatterID   int64  `json:"batterId"`
	PitcherID  int64  `json:"pitcherId"`
	EventType  string `json:"eventType"`
	Runs       int    `json:"runs"`
	RBI        int    `json:"rbi"`
	SwStrCount int
	PitchCount int
}
