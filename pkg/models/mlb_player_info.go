package models


type MLBPlayerInfo struct {
	ID               int
	FullName         string
	FirstName        string
	LastName         string
	PrimaryNumber    string
	BirthDate        string
	CurrentAge       int
	BirthCity        string
	BirthState       string
	BirthCountry     string
	Height           string
	Weight           int
	Active           bool
	BatSide          string
	PitchHand        string
	MLBDebutDate     string
	StrikeZoneTop    float32
	StrikeZoneBottom float32
	TeamID           int
	TeamName         string
	Position         string
  Season           string
}

