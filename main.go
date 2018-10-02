package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/models"
)

var users []models.Profile
var nextID = 0
var sessions = make(models.Session)

func cleanLoginInfo(r *http.Request, u *models.UserPassword) error {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &u)
	if err != nil {
		return fmt.Errorf("json error")
	}

	return nil
}

func cleanProfile(r *http.Request, p *models.Profile) error {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &p)
	if err != nil {
		return fmt.Errorf("json error")
	}

	return nil
}

func generateSessionID() (id string, err error) {
	b := make([]byte, 32)
	if _, err = io.ReadFull(rand.Reader, b); err != nil {
		return
	}
	id = base64.URLEncoding.EncodeToString(b)
	return
}

func loginUser(w http.ResponseWriter, userID int) error {
	sessionID, err := generateSessionID()
	if err != nil {
		log.Println(err)
		return err
	}
	// TODO: to db!
	sessions[sessionID] = userID

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

func sessionHandler(w http.ResponseWriter, r *http.Request) {
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

		// TODO: go to database, ask for user...

		dbResponse := models.UserPassword{}

		for _, v := range users {
			if u.Email == v.Email {
				dbResponse = v.UserPassword
			}
		}

		if u.Email == dbResponse.Email && u.Password == dbResponse.Password {
			err := loginUser(w, dbResponse.UserID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			log.Println("User logged in:", u.Email)
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
	case http.MethodDelete:
		session, err := r.Cookie("session_id")
		if err == http.ErrNoCookie {
			// user has already logged out
			return
		}

		// TODO: db
		delete(sessions, session.Value)

		session.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, session)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func validateNickname(s string) []models.ProfileError {
	var errors []models.ProfileError

	if utf8.RuneCountInString(s) < 4 {
		errors = append(errors, models.ProfileError{
			Field: "nickname",
			Text:  "Никнейм должен быть не менее 4 символов",
		})
	}
	if utf8.RuneCountInString(s) > 32 {
		errors = append(errors, models.ProfileError{
			Field: "nickname",
			Text:  "Никнейм должен быть не более 32 символов",
		})
	}
	for _, v := range users {
		if v.Nickname == s {
			errors = append(errors, models.ProfileError{
				Field: "nickname",
				Text:  "Этот никнейм уже занят",
			})
			break
		}
	}

	return errors
}

func validateEmail(s string) []models.ProfileError {
	var errors []models.ProfileError

	if strings.Contains(s, "@") == false {
		errors = append(errors, models.ProfileError{
			Field: "email",
			Text:  "Неверный формат почты",
		})
	}
	for _, v := range users {
		if v.Email == s {
			errors = append(errors, models.ProfileError{
				Field: "email",
				Text:  "Данная почта уже занята",
			})
			break
		}
	}

	return errors
}

func validatePassword(s string) []models.ProfileError {
	var errors []models.ProfileError

	if utf8.RuneCountInString(s) < 8 {
		errors = append(errors, models.ProfileError{
			Field: "password",
			Text:  "Пароль должен быть не менее 8 символов",
		})
	}
	if utf8.RuneCountInString(s) > 32 {
		errors = append(errors, models.ProfileError{
			Field: "password",
			Text:  "Пароль должен быть не более 32 символов",
		})
	}

	return errors
}

func validateFields(u *models.Profile) []models.ProfileError {
	var errors []models.ProfileError

	errors = append(errors, validateNickname(u.Nickname)...)
	errors = append(errors, validateEmail(u.Email)...)
	errors = append(errors, validatePassword(u.Email)...)

	return errors
}

func getUserIDFromSessionID(r *http.Request) (int, error) {
	c, err := r.Cookie("session_id")
	if err != nil {
		return -1, err
	}
	id := sessions[c.Value]

	return id, nil
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		q := r.URL.Query()
		id, idOK := q["id"]
		nickname, nicknameOK := q["nickname"]
		publicProfile := &models.Profile{}

		if idOK {
			intID, err := strconv.Atoi(id[0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			for _, v := range users {
				if v.UserID == intID {
					*publicProfile = v
					publicProfile.Password = ""
					w.Header().Set("Content-Type", "application/json")
					json, err := json.Marshal(publicProfile)
					if err != nil {
						log.Println(err, "in profileMethod")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					fmt.Println(string(json))
					fmt.Fprintln(w, string(json))
					return
				}
			}
		} else if nicknameOK {
			searchedNickname := nickname[0]
			for _, v := range users {
				if v.Nickname == searchedNickname {
					*publicProfile = v
					publicProfile.Password = ""
					w.Header().Set("Content-Type", "application/json")
					json, err := json.Marshal(publicProfile)
					if err != nil {
						log.Println(err, "in profileMethod")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					fmt.Println(string(json))
					fmt.Fprintln(w, string(json))
					return
				}
			}
		} else {
			searchID, err := getUserIDFromSessionID(r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			for _, v := range users {
				if searchID == v.UserID {
					*publicProfile = v
					publicProfile.Password = ""

					w.Header().Set("Content-Type", "application/json")
					json, err := json.Marshal(publicProfile)
					if err != nil {
						log.Println(err, "in profileMethod")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					fmt.Println(string(json))
					fmt.Fprintln(w, string(json))
					return
				}
			}
		}

		if publicProfile.Email == "" {
			w.WriteHeader(http.StatusNotFound)
		}

	case http.MethodPost:
		u := &models.Profile{}
		err := cleanProfile(r, u)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if u.Nickname == "" || u.Email == "" || u.Password == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		fieldErrors := validateFields(u)

		if len(fieldErrors) != 0 {
			errorsList := make(map[string][]models.ProfileError)
			errorsList["error"] = fieldErrors
			json, err := json.Marshal(errorsList)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintln(w, string(json))
		} else {
			u.UserID = nextID
			nextID++
			users = append(users, *u)

			err := loginUser(w, u.UserID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			log.Println("User logged in:", u.Email)
		}

	case http.MethodPut:
		id, err := getUserIDFromSessionID(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
		}

		u := &models.Profile{}
		err = cleanProfile(r, u)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var fieldErrors []models.ProfileError

		if u.Nickname != "" {
			fieldErrors = append(fieldErrors, validateNickname(u.Nickname)...)
		}
		if u.Email != "" {
			fieldErrors = append(fieldErrors, validateNickname(u.Email)...)
		}
		if u.Password != "" {
			fieldErrors = append(fieldErrors, validateNickname(u.Password)...)
		}

		if len(fieldErrors) != 0 {
			errorsList := make(map[string][]models.ProfileError)
			errorsList["error"] = fieldErrors
			json, err := json.Marshal(errorsList)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintln(w, string(json))
		} else {
			for k := range users {
				if id == users[k].UserID {
					users[k] = *u
					return
				}
			}
		}
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func scoreboardHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// TODO: LIMIT, OFFSET from Query
		// q := r.URL.Query()
		// limit, limitOK := q["limit"]
		// limitValue := 0
		// if limitOK {
		// 	var err error
		// 	limitValue, err = strconv.Atoi(limit[0])
		// 	if err != nil {
		// 		w.WriteHeader(http.StatusBadRequest)
		// 		return
		// 	}
		// }

		// offset, offsetOK := q["offset"]
		// offsetValue := 0
		// if offsetOK {
		// 	var err error
		// 	offsetValue, err = strconv.Atoi(offset[0])
		// 	if err != nil {
		// 		w.WriteHeader(http.StatusBadRequest)
		// 		return
		// 	}
		// }

		// TODO: db request
		var records []models.Position
		for _, v := range users {
			records = append(records, models.Position{
				ID:       v.UserID,
				Nickname: v.Nickname,
				Points:   v.Record,
			})
		}

		sort.Slice(records, func(i, j int) bool {
			return records[i].Points > records[j].Points
		})

		json, err := json.Marshal(records)
		if err != nil {
			log.Println(err, "in scoreboardHandler while parsing struct in json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(json))

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func avatarHandler(w http.ResponseWriter, r *http.Request) {
	return
}

func middlewareCORS(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Origin", "https://dmstudio.now.sh")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, User-Agent, Cache-Control, Accept, X-Requested-With, If-Modified-Since")
		next.ServeHTTP(w, r)
	})
}

func main() {
	http.HandleFunc("/session", middlewareCORS(sessionHandler))
	http.HandleFunc("/profile", middlewareCORS(profileHandler))
	http.HandleFunc("/profile/avatar", middlewareCORS(avatarHandler))
	http.HandleFunc("/scoreboard", middlewareCORS(scoreboardHandler))

	log.Println("starting server at:", 8080)
	http.ListenAndServe(":8080", nil)
}
