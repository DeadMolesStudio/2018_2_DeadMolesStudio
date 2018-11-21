package models

type Profile struct {
	User
	Nickname string  `json:"nickname" example:"Nick"`
	Avatar   *string `json:"avatar,omitempty"`
	Stats
}

type RegisterProfile struct {
	Nickname string `json:"nickname" example:"Nick"`
	UserPassword
}

type User struct {
	UserID uint `json:"id" db:"user_id"`
	UserPassword
}

type UserPassword struct {
	Email    string `json:"email" example:"email@email.com" valid:"required~Почта не может быть пустой,email~Невалидная почта"`
	Password string `json:"password,omitempty" example:"password" valid:"stringlength(8|32)~Пароль должен быть не менее 8 символов и не более 32 символов"`
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
	ID       uint
	Nickname string
}
