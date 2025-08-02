package models

type PitcherEvent struct {
	PitcherID int64  `json:"pitcherId"`
	EventType string `json:"eventType"` // e.g. "strikeout", "walk", etc.
	// SwStrCount int64
	// PitchCount int
}
