package models

type Position struct {
	ID       uint   `json:"id" example:"42" db:"user_id"`
	Nickname string `json:"nickname" example:"Nick"`
	Points   int    `json:"record" example:"100500" db:"record"`
}

type PositionList struct {
	List  []Position `json:"players"`
	Total int        `json:"total" example:"1"`
}

type FetchScoreboardPage struct {
	Limit uint
	Page  uint
}
