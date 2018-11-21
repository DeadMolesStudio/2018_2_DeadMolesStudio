package main

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/logger"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/middleware"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/session"

	"api/database"
	_ "api/docs"
	"api/filesystem"
	"api/handlers"
)

func main() {
	l := logger.InitLogger()
	defer l.Sync()

	db := database.InitDB("postgres@postgres:5432", "ketnipz")
	defer db.Close()

	sm := session.ConnectSessionManager()
	defer sm.Close()

	http.HandleFunc("/session", middleware.RecoverMiddleware(middleware.AccessLogMiddleware(
		middleware.CORSMiddleware(middleware.SessionMiddleware(handlers.SessionHandler)))))
	http.HandleFunc("/profile", middleware.RecoverMiddleware(middleware.AccessLogMiddleware(
		middleware.CORSMiddleware(middleware.SessionMiddleware(handlers.ProfileHandler)))))
	http.HandleFunc("/profile/avatar", middleware.RecoverMiddleware(middleware.AccessLogMiddleware(
		middleware.CORSMiddleware(middleware.SessionMiddleware(handlers.AvatarHandler)))))
	http.HandleFunc("/scoreboard", middleware.RecoverMiddleware(middleware.AccessLogMiddleware(
		middleware.CORSMiddleware(handlers.ScoreboardHandler))))

	// swag init -g handlers/api.go
	http.HandleFunc("/api/docs/", httpSwagger.WrapHandler)

	http.HandleFunc("/static/", middleware.RecoverMiddleware(middleware.AccessLogMiddleware(
		middleware.CORSMiddleware(filesystem.StaticHandler))))

	logger.Info("starting server at: ", 8080)
	logger.Panic(http.ListenAndServe(":8080", nil))
}
