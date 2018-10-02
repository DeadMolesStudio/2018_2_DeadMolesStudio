package main

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/handlers"
)

func middlewareCORS(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Origin", "https://dmstudio.now.sh")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, User-Agent, Cache-Control, Accept, X-Requested-With, If-Modified-Since")
		next.ServeHTTP(w, r)
	})
}

func main() {
	http.HandleFunc("/session", middlewareCORS(handlers.SessionHandler))
	http.HandleFunc("/profile", middlewareCORS(handlers.ProfileHandler))
	http.HandleFunc("/profile/avatar", middlewareCORS(handlers.AvatarHandler))
	http.HandleFunc("/scoreboard", middlewareCORS(handlers.ScoreboardHandler))

	log.Println("starting server at:", 8080)
	http.ListenAndServe(":8080", nil)
}
