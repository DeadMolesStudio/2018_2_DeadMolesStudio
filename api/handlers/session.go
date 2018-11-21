package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/logger"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/middleware"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/session"

	"api/database"
	"api/models"
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

func loginUser(w http.ResponseWriter, userID uint) error {
	sessionID, err := session.Create(userID)
	if err != nil {
		logger.Error(err)
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

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getSession(w, r)
	case http.MethodPost:
		postSession(w, r)
	case http.MethodDelete:
		deleteSession(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// @Summary Получить сессию
// @Description Получить сессию пользователя, если есть сессия, то она в куке session_id
// @ID get-session
// @Produce json
// @Success 200 {object} models.Session "Пользователь залогинен, успешно"
// @Failure 401 "Не залогинен"
// @Failure 500 "Ошибка в бд"
// @Router /session [GET]
func getSession(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value(middleware.KeyIsAuthenticated).(bool) {
		sID, err := json.Marshal(models.Session{SessionID: r.Context().Value(middleware.KeySessionID).(string)})
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(sID))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

// @Summary Залогинить
// @Description Залогинить пользователя (создать сессию)
// @ID post-session
// @Accept json
// @Produce json
// @Param UserPassword body models.UserPassword true "Почта и пароль"
// @Success 200 {object} models.Session "Успешный вход / пользователь уже залогинен"
// @Failure 400 "Неверный формат JSON, невалидные данные"
// @Failure 422 "Неверная пара пользователь/пароль"
// @Failure 500 "Внутренняя ошибка"
// @Router /session [POST]
func postSession(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value(middleware.KeyIsAuthenticated).(bool) {
		// user has already logged in
		return
	}

	u := &models.UserPassword{}
	err := cleanLoginInfo(r, u)
	if err != nil {
		switch err.(type) {
		case ParseJSONError:
			w.WriteHeader(http.StatusBadRequest)
		default:
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	isValid := govalidator.IsEmail(u.Email)
	if !isValid {
		sendError(w, r, fmt.Errorf("Невалидная почта"), http.StatusBadRequest)
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
		logger.Info("user with id %v and email %v logged in", dbResponse.UserID, dbResponse.Email)
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
}

// @Summary Разлогинить
// @Descriptiond Разлогинить пользователя (удалить сессию)
// @ID delete-session
// @Success 200 "Успешный выход / пользователь уже разлогинен"
// @Router /session [DELETE]
func deleteSession(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value(middleware.KeyIsAuthenticated).(bool) {
		// user has already logged out
		return
	}
	err := session.Delete(r.Context().Value(middleware.KeySessionID).(string))
	if err != nil { // but we continue
		logger.Error(err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Expires:  time.Now().AddDate(0, 0, -1),
		Secure:   true,
		HttpOnly: true,
	})
}
