package models

type BattingOrder struct {
	TeamAbbreviation string `csv:"teamAbbreviation"`
	PitcherThrows    string `csv:"pitcher_throws"`
	PlayerID1        int    `csv:"playerId_1"`
	PlayerID2        int    `csv:"playerId_2"`
	PlayerID3        int    `csv:"playerId_3"`
	PlayerID4        int    `csv:"playerId_4"`
	PlayerID5        int    `csv:"playerId_5"`
	PlayerID6        int    `csv:"playerId_6"`
	PlayerID7        int    `csv:"playerId_7"`
	PlayerID8        int    `csv:"playerId_8"`
	PlayerID9        int    `csv:"playerId_9"`
	N1               int    `csv:"n_1"`
	N2               int    `csv:"n_2"`
	N3               int    `csv:"n_3"`
	N4               int    `csv:"n_4"`
	N5               int    `csv:"n_5"`
	N6               int    `csv:"n_6"`
	N7               int    `csv:"n_7"`
	N8               int    `csv:"n_8"`
	N9               int    `csv:"n_9"`
}
