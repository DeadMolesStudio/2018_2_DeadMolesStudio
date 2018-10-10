package models

type Profile struct {
	Nickname string `json:"nickname"`
	UserPassword
	Stats
}

type UserPassword struct {
	UserID   int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type Stats struct {
	Record int `json:"record"`
	Win    int `json:"win"`
	Draws  int `json:"draws"`
	Loss   int `json:"loss"`
}

type ProfileError struct {
	Field string `json:"field"`
	Text  string `json:"text"`
}

type ProfileErrorList struct {
	Errors []ProfileError `json:"error"`
}

type RequestProfile struct {
	ID       int
	Nickname string
}
