package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/satori/go.uuid"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/database"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/models"
)

func cleanLoginInfo(r *http.Request, u *models.UserPassword) error {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, u)
	if err != nil {
		return ParseJSONError{err}
	}

	return nil
}

func loginUser(w http.ResponseWriter, userID int) error {
	sessionID := ""
	for {
		var err error
		sessionID = uuid.NewV4().String()
		exists, err := database.CheckExistenceOfSession(sessionID)
		if err != nil {
			log.Println(err)
			return err
		}
		if !exists {
			break
		}
	}

	err := database.CreateNewSession(sessionID, userID)
	if err != nil {
		log.Println(err)
		return err
	}

	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	return nil
}

func getUserIDFromSessionID(r *http.Request) (int, error) {
	c, err := r.Cookie("session_id")
	if err != nil {
		return -1, err
	}
	id, err := database.GetIDFromSession(c.Value)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getSession(w, r)
	case http.MethodPost:
		postSession(w, r)
	case http.MethodDelete:
		deleteSession(w, r)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// @Title Получить сессию
// @Summary Получить сессию пользователя, если есть сессия, то она в куке session_id
// @ID get-session
// @Produce json
// @Success 200 {object} models.Session "Пользователь залогинен, успешно"
// @Failure 401 "Не залогинен"
// @Failure 500 "Ошибка в бд"
// @Router /session [GET]
func getSession(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_id")
	if err == nil {
		sID, err := json.Marshal(models.Session{SessionID: c.Value})
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(sID))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

// @Title Залогинить
// @Summary Залогинить пользователя (создать сессию)
// @ID post-session
// @Accept json
// @Produce json
// @Param UserPassword body models.UserPassword true "Почта и пароль"
// @Success 200 {object} models.Session "Успешный вход / пользователь уже залогинен"
// @Failure 400 "Неверный формат JSON"
// @Failure 422 "Неверная пара пользователь/пароль"
// @Failure 500 "Внутренняя ошибка"
// @Router /session [POST]
func postSession(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_id")
	if err == nil {
		// user has already logged in
		return
	}

	u := &models.UserPassword{}
	err = cleanLoginInfo(r, u)
	if err != nil {
		switch err.(type) {
		case ParseJSONError:
			w.WriteHeader(http.StatusBadRequest)
		default:
			log.Println(err, "in sessionHandler in getUserFromRequestBody")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if u.Email == "" || u.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbResponse, err := database.GetUserPassword(u.Email)

	if err != nil {
		switch err.(type) {
		case database.UserNotFoundError:
			w.WriteHeader(http.StatusUnprocessableEntity)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if u.Email == dbResponse.Email && u.Password == dbResponse.Password {
		err := loginUser(w, dbResponse.UserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println("User logged in:", dbResponse.UserID, dbResponse.Email)
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
}

// @Title Разлогинить
// @Summary Разлогинить пользователя (удалить сессию)
// @ID delete-session
// @Success 200 "Успешный выход / пользователь уже разлогинен"
// @Router /session [DELETE]
func deleteSession(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		// user has already logged out
		return
	}

	err = database.DeleteSession(session.Value)
	if err != nil { // but we continue
		log.Println(err)
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
}
