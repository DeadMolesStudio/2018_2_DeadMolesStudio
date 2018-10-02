package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/database"
)

func ScoreboardHandler(w http.ResponseWriter, r *http.Request) {
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

		records, err := database.GetUserPositionsDescending()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

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
