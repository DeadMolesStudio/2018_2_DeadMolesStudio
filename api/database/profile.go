package database

import (
	"database/sql"
	"strings"

	"github.com/lib/pq"

	"api/models"
)

func GetUserPassword(e string) (models.User, error) {
	res := models.User{}
	qres := db.QueryRowx(`
		SELECT user_id, email, password FROM user_profile
		WHERE email = $1`,
		e)
	if err := qres.Err(); err != nil {
		return res, err
	}
	err := qres.StructScan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, UserNotFoundError{"email"}
		}
		return res, err
	}

	return res, nil
}

func CreateNewUser(u *models.RegisterProfile) (models.Profile, error) {
	res := models.Profile{}
	qres := db.QueryRowx(`
		INSERT INTO user_profile (email, password, nickname)
		VALUES ($1, $2, $3) RETURNING user_id, email, nickname`,
		u.Email, u.Password, u.Nickname)
	if err := qres.Err(); err != nil {
		pqErr := err.(*pq.Error)
		switch pqErr.Code {
		case "23502":
			return res, ErrNotNullConstraintViolation
		case "23505":
			return res, ErrUniqueConstraintViolation
		}
	}
	err := qres.StructScan(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func UpdateUserByID(id uint, u *models.RegisterProfile) error {
	if u.Email == "" && u.Password == "" && u.Nickname == "" {
		return nil
	}

	q := strings.Builder{}
	q.WriteString(`
		UPDATE user_profile
		SET `)
	hasBefore := false
	if u.Email != "" {
		q.WriteString("email = :email")
		hasBefore = true
	}
	if u.Password != "" {
		if hasBefore {
			q.WriteString(", ")
		}
		q.WriteString("password = :password")
		hasBefore = true
	}
	if u.Nickname != "" {
		if hasBefore {
			q.WriteString(", ")
		}
		q.WriteString("nickname = :nickname")
	}
	q.WriteString(`
		WHERE user_id = :user_id`)

	_, err := db.NamedExec(q.String(), &models.Profile{
		User: models.User{
			UserID: id,
			UserPassword: models.UserPassword{
				Email:    u.Email,
				Password: u.Password,
			},
		},
		Nickname: u.Nickname,
	})
	if err != nil {
		return err
	}

	return nil
}

func GetUserProfileByID(id uint) (models.Profile, error) {
	res := models.Profile{}
	qres := db.QueryRowx(`
		SELECT user_id, email, nickname, avatar, record, win, draws, loss FROM user_profile
		WHERE user_id = $1`,
		id)
	if err := qres.Err(); err != nil {
		return res, err
	}
	err := qres.StructScan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, UserNotFoundError{"id"}
		}
		return res, err
	}

	return res, nil
}

func GetUserProfileByNickname(nickname string) (models.Profile, error) {
	res := models.Profile{}
	qres := db.QueryRowx(`
		SELECT user_id, email, nickname, avatar, record, win, draws, loss FROM user_profile
		WHERE nickname = $1`,
		nickname)
	if err := qres.Err(); err != nil {
		return res, err
	}
	err := qres.StructScan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, UserNotFoundError{"nickname"}
		}
		return res, err
	}

	return res, nil
}

func CheckExistenceOfEmail(e string) (bool, error) {
	res := models.Profile{}
	qres := db.QueryRowx(`
		SELECT FROM user_profile
		WHERE email = $1`,
		e)
	if err := qres.Err(); err != nil {
		return false, err
	}
	err := qres.StructScan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func CheckExistenceOfNickname(n string) (bool, error) {
	res := models.Profile{}
	qres := db.QueryRowx(`
		SELECT FROM user_profile
		WHERE nickname = $1`,
		n)
	if err := qres.Err(); err != nil {
		return false, err
	}
	err := qres.StructScan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func GetCountOfUsers() (int, error) {
	res := 0
	// TODO: optimize it
	qres := db.QueryRowx(`
		SELECT COUNT(*) FROM user_profile`)
	if err := qres.Err(); err != nil {
		return 0, err
	}
	err := qres.Scan(&res)
	if err != nil {
		return 0, err
	}

	return res, nil
}

func UploadAvatar(uID uint, path string) error {
	qres, err := db.Exec(`
		UPDATE user_profile
		SET avatar = $2
		WHERE user_id = $1`,
		uID, path)
	if err != nil {
		return err
	}
	res, err := qres.RowsAffected()
	if err != nil {
		return err
	}
	if res == 0 {
		return &UserNotFoundError{"id"}
	}

	return nil
}

func DeleteAvatar(uID uint) error {
	qres, err := db.Exec(`
		UPDATE user_profile
		SET avatar = NULL
		WHERE user_id = $1`,
		uID)
	if err != nil {
		return err
	}
	res, err := qres.RowsAffected()
	if err != nil {
		return err
	}
	if res == 0 {
		return &UserNotFoundError{"id"}
	}

	return nil
}
