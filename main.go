package main

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/database"
	_ "github.com/go-park-mail-ru/2018_2_DeadMolesStudio/docs"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/handlers"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/logger"
)

func main() {
	l := logger.InitLogger()
	defer l.Sync()

	db := database.InitDB("postgres@postgres:5432", "ketnipz")
	defer db.Close()

	http.HandleFunc("/session", handlers.RecoverMiddleware(handlers.AccessLogMiddleware(
		handlers.CORSMiddleware(handlers.SessionMiddleware(handlers.SessionHandler)))))
	http.HandleFunc("/profile", handlers.RecoverMiddleware(handlers.AccessLogMiddleware(
		handlers.CORSMiddleware(handlers.SessionMiddleware(handlers.ProfileHandler)))))
	http.HandleFunc("/profile/avatar", handlers.RecoverMiddleware(handlers.AccessLogMiddleware(
		handlers.CORSMiddleware(handlers.SessionMiddleware(handlers.AvatarHandler)))))
	http.HandleFunc("/scoreboard", handlers.RecoverMiddleware(handlers.AccessLogMiddleware(
		handlers.CORSMiddleware(handlers.ScoreboardHandler))))

	// swag init -g handlers/api.go
	http.HandleFunc("/api/docs/", httpSwagger.WrapHandler)

	logger.Info("starting server at: ", 8080)
	logger.Panic(http.ListenAndServe(":8080", nil))
}
