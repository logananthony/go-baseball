package models

type BatterEvent struct {
	BatterID  int64  `json:"batterId"`
	EventType string `json:"eventType"` // e.g. "single", "double", "home_run", etc.
}
