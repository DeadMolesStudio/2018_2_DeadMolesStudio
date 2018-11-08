package filesystem

import (
	"net/http"
)

var fs = http.StripPrefix("/static/", http.FileServer(http.Dir("static")))

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	fs.ServeHTTP(w, r)
}
