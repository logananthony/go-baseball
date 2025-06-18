package models

import "time"

type GameDataGamePk struct {
	GamePk           int
	Season           int
	GameDate         time.Time
	Team             string
	TeamAbbreviation string
	BattingOrder     int
	PlayerId         int
	PlayerName       string
	HomePitcherId    int
	HomePitcherName  string
	AwayPitcherId    int
	AwayPitcherName  string
	HomeTeamAbbr     string
	AwayTeamAbbr     string
}
