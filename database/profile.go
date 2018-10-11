package database

import (
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/models"
)

var users []models.Profile
var nextID = 1

func GetUserPassword(e string) (models.UserPassword, error) {
	for _, v := range users {
		if e == v.Email {
			return v.UserPassword, nil
		}
	}

	return models.UserPassword{}, UserNotFoundError{"email"}
}

func CreateNewUser(u *models.Profile) error {
	u.UserID = nextID
	nextID++
	users = append(users, *u)

	return nil
}

func UpdateUserByID(id int, u *models.Profile) error {
	for k := range users {
		if id == users[k].UserID {
			if u.Nickname != "" {
				users[k].Nickname = u.Nickname
			}
			if u.Email != "" {
				users[k].Email = u.Email
			}
			if u.Password != "" {
				users[k].Password = u.Password
			}
			return nil
		}
	}

	return UserNotFoundError{"id"}
}

func GetUserProfileByID(id int) (models.Profile, error) {
	for _, v := range users {
		if id == v.UserID {
			v.Password = ""
			return v, nil
		}
	}

	return models.Profile{}, UserNotFoundError{"id"}
}

func GetUserProfileByNickname(nickname string) (models.Profile, error) {
	for _, v := range users {
		if nickname == v.Nickname {
			v.Password = ""
			return v, nil
		}
	}

	return models.Profile{}, UserNotFoundError{"nickname"}
}

func CheckExistenceOfEmail(e string) (bool, error) {
	for _, v := range users {
		if e == v.Email {
			return true, nil
		}
	}

	return false, nil
}

func CheckExistenceOfNickname(n string) (bool, error) {
	for _, v := range users {
		if n == v.Nickname {
			return true, nil
		}
	}

	return false, nil
}
