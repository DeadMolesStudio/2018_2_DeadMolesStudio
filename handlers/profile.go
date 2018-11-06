package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/database"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/models"
)

func cleanProfile(r *http.Request, p *models.RegisterProfile) error {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, p)
	if err != nil {
		return ParseJSONError{err}
	}

	return nil
}

func validateNickname(s string) ([]models.ProfileError, error) {
	var errors []models.ProfileError

	isValid := govalidator.StringLength(s, "4", "32")
	if !isValid {
		errors = append(errors, models.ProfileError{
			Field: "nickname",
			Text:  "Никнейм должен быть не менее 4 символов и не более 32 символов",
		})
		return errors, nil
	}

	exists, err := database.CheckExistenceOfNickname(s)
	if err != nil {
		log.Println(err)
		return errors, err
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

	isValid := govalidator.IsEmail(s)
	if !isValid {
		errors = append(errors, models.ProfileError{
			Field: "email",
			Text:  "Невалидная почта",
		})
		return errors, nil
	}

	exists, err := database.CheckExistenceOfEmail(s)
	if err != nil {
		log.Println(err)
		return errors, err
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

	isValid := govalidator.StringLength(s, "8", "32")
	if !isValid {
		errors = append(errors, models.ProfileError{
			Field: "password",
			Text:  "Пароль должен быть не менее 8 символов и не более 32 символов",
		})
	}

	return errors
}

func validateFields(u *models.RegisterProfile) ([]models.ProfileError, error) {
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
	errors = append(errors, validatePassword(u.Password)...)

	return errors, nil
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getProfile(w, r)
	case http.MethodPost:
		postProfile(w, r)
	case http.MethodPut:
		putProfile(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// @Title Получить профиль
// @Summary Получить профиль пользователя по ID, email или из сессии
// @ID get-profile
// @Produce json
// @Param id query uint false "ID"
// @Param nickname query string false "Никнейм"
// @Success 200 {object} models.Profile "Пользователь найден, успешно"
// @Failure 400 "Неправильный запрос"
// @Failure 401 "Не залогинен"
// @Failure 404 "Не найдено"
// @Failure 500 "Ошибка в бд"
// @Router /profile [GET]
func getProfile(w http.ResponseWriter, r *http.Request) {
	params := &models.RequestProfile{}
	err := decoder.Decode(params, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if params.ID != 0 {
		profile, err := database.GetUserProfileByID(params.ID)
		if err != nil {
			switch err.(type) {
			case database.UserNotFoundError:
				w.WriteHeader(http.StatusNotFound)
				return
			default:
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
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
			switch err.(type) {
			case database.UserNotFoundError:
				w.WriteHeader(http.StatusNotFound)
				return
			default:
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
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
		if !r.Context().Value(keyIsAuthenticated).(bool) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		profile, err := database.GetUserProfileByID(r.Context().Value(keyUserID).(uint))
		if err != nil {
			switch err.(type) {
			case database.UserNotFoundError:
				w.WriteHeader(http.StatusNotFound)
				return
			default:
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
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
}

// @Title Зарегистрироваться и залогиниться по новому профилю
// @Summary Зарегистрировать по никнейму, почте и паролю и автоматически залогинить
// @ID post-profile
// @Accept json
// @Produce json
// @Param Profile body models.RegisterProfile true "Никнейм, почта и пароль"
// @Success 200 "Пользователь зарегистрирован и залогинен успешно"
// @Failure 400 "Неверный формат JSON"
// @Failure 403 {object} models.ProfileErrorList "Ошибки при регистрации: невалидна или занята почта, занят ник, пароль не удовлетворяет правилам безопасности, другие ошибки"
// @Failure 422 "При регистрации не все параметры"
// @Failure 500 "Ошибка в бд"
// @Router /profile [POST]
func postProfile(w http.ResponseWriter, r *http.Request) {
	u := &models.RegisterProfile{}
	err := cleanProfile(r, u)
	if err != nil {
		switch err.(type) {
		case ParseJSONError:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
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
		json, err := json.Marshal(models.ProfileErrorList{Errors: fieldErrors})
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, string(json))
	} else {
		newU, err := database.CreateNewUser(u)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = loginUser(w, newU.UserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println("New user logged in:", newU.UserID, newU.Email, newU.Nickname)
	}
}

// @Title Изменить профиль
// @Summary Изменить профиль, должен быть залогинен
// @ID put-profile
// @Accept json
// @Produce json
// @Param Profile body models.RegisterProfile true "Новые никнейм, и/или почта, и/или пароль"
// @Success 200 "Пользователь найден, успешно изменены данные"
// @Failure 400 "Неверный формат JSON"
// @Failure 401 "Не залогинен"
// @Failure 403 {object} models.ProfileErrorList "Ошибки при регистрации: невалидна или занята почта, занят ник, пароль не удовлетворяет правилам безопасности, другие ошибки"
// @Failure 500 "Ошибка в бд"
// @Router /profile [PUT]
func putProfile(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value(keyIsAuthenticated).(bool) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	u := &models.RegisterProfile{}
	err := cleanProfile(r, u)
	if err != nil {
		switch err.(type) {
		case ParseJSONError:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if u.Nickname == "" && u.Email == "" && u.Password == "" {
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
		json, err := json.Marshal(models.ProfileErrorList{Errors: fieldErrors})
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, string(json))
	} else {
		id := r.Context().Value(keyUserID).(uint)
		err := database.UpdateUserByID(id, u)
		if err != nil {
			switch err.(type) {
			case database.UserNotFoundError:
				w.WriteHeader(http.StatusNotFound)
			default:
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		log.Println("User with id", id, "changed to", u.Nickname, u.Email)
	}
}

func AvatarHandler(w http.ResponseWriter, r *http.Request) {
	return
}
