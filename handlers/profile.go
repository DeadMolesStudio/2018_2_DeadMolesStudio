package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"

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

	err = json.Unmarshal(body, p)
	if err != nil {
		return fmt.Errorf("json error")
	}

	return nil
}

func validateNickname(s string) ([]models.ProfileError, error) {
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

	exists, err := database.CheckExistenceOfNickname(s)
	if err != nil {
		log.Println(err)
		return []models.ProfileError{}, err
	}
	if exists {
		errors = append(errors, models.ProfileError{
			Field: "nickname",
			Text:  "Этот никнейм уже занят",
		})
	}

	return errors, nil
}

func validateEmail(s string) ([]models.ProfileError, error) {
	var errors []models.ProfileError

	if strings.Contains(s, "@") == false {
		errors = append(errors, models.ProfileError{
			Field: "email",
			Text:  "Неверный формат почты",
		})
	}

	exists, err := database.CheckExistenceOfEmail(s)
	if err != nil {
		log.Println(err)
		return []models.ProfileError{}, err
	}
	if exists {
		errors = append(errors, models.ProfileError{
			Field: "email",
			Text:  "Данная почта уже занята",
		})
	}

	return errors, nil
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

func validateFields(u *models.Profile) ([]models.ProfileError, error) {
	var errors []models.ProfileError

	valErrors, dbErr := validateNickname(u.Nickname)
	if dbErr != nil {
		return []models.ProfileError{}, dbErr
	}
	errors = append(errors, valErrors...)

	valErrors, dbErr = validateEmail(u.Email)
	if dbErr != nil {
		return []models.ProfileError{}, dbErr
	}
	errors = append(errors, valErrors...)
	errors = append(errors, validatePassword(u.Email)...)

	return errors, nil
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		params := &models.RequestProfile{}
		err := decoder.Decode(params, r.URL.Query())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if params.ID != 0 {
			profile, err := database.GetUserProfileByID(params.ID)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json, err := json.Marshal(profile)
			if err != nil {
				log.Println(err, "in profileMethod")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fmt.Fprintln(w, string(json))
		} else if params.Nickname != "" {
			profile, err := database.GetUserProfileByNickname(params.Nickname)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json, err := json.Marshal(profile)
			if err != nil {
				log.Println(err, "in profileMethod")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fmt.Fprintln(w, string(json))
		} else {
			searchID, err := getUserIDFromSessionID(r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			profile, err := database.GetUserProfileByID(searchID)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json, err := json.Marshal(profile)
			if err != nil {
				log.Println(err, "in profileMethod")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fmt.Fprintln(w, string(json))
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

		fieldErrors, err := validateFields(u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(fieldErrors) != 0 {
			errorsList := models.ProfileErrorList{
				Errors: fieldErrors,
			}
			json, err := json.Marshal(errorsList)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintln(w, string(json))
		} else {
			err := database.CreateNewUser(u)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = loginUser(w, u.UserID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Println("New user logged in:", u.UserID, u.Email, u.Nickname)
		}

	case http.MethodPut:
		id, err := getUserIDFromSessionID(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		u := &models.Profile{}
		err = cleanProfile(r, u)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var fieldErrors []models.ProfileError

		if u.Nickname != "" {
			valErrors, dbErr := validateNickname(u.Nickname)
			if dbErr != nil {
				log.Println(dbErr)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fieldErrors = append(fieldErrors, valErrors...)
		}
		if u.Email != "" {
			valErrors, dbErr := validateEmail(u.Email)
			if dbErr != nil {
				log.Println(dbErr)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fieldErrors = append(fieldErrors, valErrors...)
		}
		if u.Password != "" {
			fieldErrors = append(fieldErrors, validatePassword(u.Password)...)
		}

		if len(fieldErrors) != 0 {
			errorsList := make(map[string][]models.ProfileError)
			errorsList["error"] = fieldErrors
			json, err := json.Marshal(errorsList)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintln(w, string(json))
		} else {
			err := database.UpdateUserByID(id, u)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			log.Println("User with id", id, "changed to", u.Nickname, u.Email)
		}
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func AvatarHandler(w http.ResponseWriter, r *http.Request) {
	return
}
