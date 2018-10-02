package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/database"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/models"
)

func generateSessionID() (id string, err error) {
	b := make([]byte, 32)
	if _, err = io.ReadFull(rand.Reader, b); err != nil {
		return
	}
	id = base64.URLEncoding.EncodeToString(b)
	return
}

func loginUser(w http.ResponseWriter, userID int) error {
	sessionID := ""
	for {
		var err error
		sessionID, err = generateSessionID()
		if err != nil {
			log.Println(err)
			return err
		}
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
		Name:    "session_id",
		Value:   sessionID,
		Expires: time.Now().Add(30 * 24 * time.Hour),
		// Secure:  true,
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
		c, err := r.Cookie("session_id")
		if err == nil {
			sID, err := json.Marshal(map[string]string{
				c.Name: c.Value,
			})
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
	case http.MethodPost:
		_, err := r.Cookie("session_id")
		if err == nil {
			// user has already logged in
			return
		}

		u := &models.UserPassword{}
		err = cleanLoginInfo(r, u)
		if invalid := u.Email == "" || u.Password == ""; err != nil || invalid {
			if invalid || err.Error() == "json error" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			log.Println(err, "in sessionHandler in getUserFromRequestBody")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		dbResponse, err := database.GetUserPassword(u.Email)

		if err != nil { // TODO: db request fail
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		if u.Email == dbResponse.Email && u.Password == dbResponse.Password {
			err := loginUser(w, dbResponse.UserID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Println("User logged in:", u.UserID, u.Email)
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
	case http.MethodDelete:
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
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
