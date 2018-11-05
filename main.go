package main

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/handlers"
)

func main() {
	http.HandleFunc("/session", handlers.CORSMiddleware(handlers.SessionHandler))
	http.HandleFunc("/profile", handlers.CORSMiddleware(handlers.ProfileHandler))
	http.HandleFunc("/profile/avatar", handlers.CORSMiddleware(handlers.AvatarHandler))
	http.HandleFunc("/scoreboard", handlers.CORSMiddleware(handlers.ScoreboardHandler))

	log.Println("starting server at:", 8080)
	log.Panicln(http.ListenAndServe(":8080", nil))
}
