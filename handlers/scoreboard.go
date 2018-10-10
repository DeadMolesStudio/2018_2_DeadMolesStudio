package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/database"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/models"
)

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
