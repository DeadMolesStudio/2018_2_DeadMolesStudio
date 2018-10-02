package models

type Position struct {
	ID       int    `json:"id"`
	Nickname string `json:"nickname"`
	Points   int    `json:"record"`
}
