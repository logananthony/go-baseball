package models

type BullpenOrder struct {
	TeamAbbreviation string `csv:"abbreviation"`
	RosterTeamID     int    `csv:"rosterTeamId"`
	PlayerID1        int    `csv:"playerId_1"`
	PlayerID2        int    `csv:"playerId_2"`
	PlayerID3        int    `csv:"playerId_3"`
	PlayerID4        int    `csv:"playerId_4"`
	PlayerID5        int    `csv:"playerId_5"`
	PlayerID6        int    `csv:"playerId_6"`
	PlayerID7        int    `csv:"playerId_7"`
	PlayerID8        int    `csv:"playerId_8"`
}
