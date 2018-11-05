package models

type Profile struct {
	User
	Nickname string `json:"nickname" example:"Nick"`
	Stats
}

type RegisterProfile struct {
	Nickname string `json:"nickname" example:"Nick"`
	UserPassword
}

type User struct {
	UserID int `json:"id"`
	UserPassword
}

type UserPassword struct {
	Email    string `json:"email" example:"email@email.com"`
	Password string `json:"password,omitempty" example:"password"`
}

type Stats struct {
	Record int `json:"record"`
	Win    int `json:"win"`
	Draws  int `json:"draws"`
	Loss   int `json:"loss"`
}

type ProfileError struct {
	Field string `json:"field" example:"nickname"`
	Text  string `json:"text" example:"Этот никнейм уже занят"`
}

type ProfileErrorList struct {
	Errors []ProfileError `json:"error"`
}

type RequestProfile struct {
	ID       int
	Nickname string
}
