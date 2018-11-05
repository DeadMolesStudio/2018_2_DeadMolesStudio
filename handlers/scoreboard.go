package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/database"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/models"
)

// @Title Получить таблицу лидеров
// @Summary Получить таблицу лидеров (пагинация присутствует)
// @ID get-scoreboard
// @Produce json
// @Param Limit query int false "Пользователей на страницу"
// @Param Page query int false "Страница номер"
// @Success 200 {object} models.PositionList "Таблицу лидеров или ее страница и общее количество"
// @Failure 500 "Ошибка в бд"
// @Router /scoreboard [GET]
func ScoreboardHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		params := &models.FetchScoreboardPage{}
		err := decoder.Decode(params, r.URL.Query())

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		records, total, err := database.GetUserPositionsDescendingPaginated(
			params)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		positionsList := models.PositionList{
			List:  records,
			Total: total,
		}
		json, err := json.Marshal(positionsList)
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
