package main

import (
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/go-park-mail-ru/2018_2_DeadMolesStudio/docs"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/handlers"
)

func main() {
	http.HandleFunc("/session", handlers.CORSMiddleware(handlers.SessionHandler))
	http.HandleFunc("/profile", handlers.CORSMiddleware(handlers.ProfileHandler))
	http.HandleFunc("/profile/avatar", handlers.CORSMiddleware(handlers.AvatarHandler))
	http.HandleFunc("/scoreboard", handlers.CORSMiddleware(handlers.ScoreboardHandler))

	// swag init -g handlers/api.go
	http.HandleFunc("/api/docs/", httpSwagger.WrapHandler)

	log.Println("starting server at:", 8080)
	log.Panicln(http.ListenAndServe(":8080", nil))
}
