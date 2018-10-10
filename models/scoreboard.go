package models

type Position struct {
	ID       int    `json:"id"`
	Nickname string `json:"nickname"`
	Points   int    `json:"record"`
}

type PositionList struct {
	List  []Position `json:"players"`
	Total int        `json:"total"`
}

type FetchScoreboardPage struct {
	Limit int
	Page  int
}
