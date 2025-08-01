package models

type BatterEvent struct {
	BatterID  int64  `json:"batterId"`
	EventType string `json:"eventType"`
	Runs      int    `json:"runs"`
	RBI       int    `json:"rbi"`
}
