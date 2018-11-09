package handlers

import (
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/game"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/logger"
)

// @Summary Начать игру по WebSocket
// @Description Инициализирует соединение для пользователя
// @ID get-game-ws
// @Success 101 "Switching Protocols"
// @Failure 400 "Нет нужных заголовков"
// @Router /game/ws [GET]
func StartGame(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Cannot upgrade connection: ", err)
		return
	}

	game.AddPlayer(conn)
}
