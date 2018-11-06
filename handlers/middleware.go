package handlers

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/database"
)

type key int

const (
	keyIsAuthenticated key = iota
	keySessionID
	keyUserID
)

func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Origin", "https://dmstudio.now.sh")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, User-Agent, Cache-Control, Accept, X-Requested-With, If-Modified-Since, Origin")

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SessionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		c, err := r.Cookie("session_id")
		if err == nil {
			uid, err := database.GetIDFromSession(c.Value)
			switch err {
			case nil:
				ctx = context.WithValue(ctx, keyIsAuthenticated, true)
				ctx = context.WithValue(ctx, keySessionID, c.Value)
				ctx = context.WithValue(ctx, keyUserID, uid)
			case database.ErrSessionNotFound:
				// delete unvalid cookie
				c.Expires = time.Now().AddDate(0, 0, -1)
				http.SetCookie(w, c)
				ctx = context.WithValue(ctx, keyIsAuthenticated, false)
			default:
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else { // ErrNoCookie
			ctx = context.WithValue(ctx, keyIsAuthenticated, false)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RecoverMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("[PANIC]:", err, "at", string(debug.Stack()))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
