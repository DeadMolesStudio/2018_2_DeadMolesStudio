package models

type Position struct {
	ID       int    `json:"id" example:"42"`
	Nickname string `json:"nickname" example:"Nick"`
	Points   int    `json:"record" example:"100500"`
}

type PositionList struct {
	List  []Position `json:"players"`
	Total int        `json:"total" example:"1"`
}

type FetchScoreboardPage struct {
	Limit int
	Page  int
}
