package models

type MLBTeamId struct {
	TeamID       int    `csv:"teamId"`
	Code         string `csv:"code"`
	FileCode     string `csv:"fileCode"`
	Abbreviation string `csv:"abbreviation"`
	Name         string `csv:"name"`
	FullName     string `csv:"fullName"`
	BriefName    string `csv:"briefName"`
}
