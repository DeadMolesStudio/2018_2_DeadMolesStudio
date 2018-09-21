package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type User struct {
	User     string
	Password string
}

func cleanUser(r *http.Request, u *User) error {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, u)
	if err != nil {
		return fmt.Errorf("json error")
	}

	return nil
}

func generateSessionID() (id string, err error) {
	b := make([]byte, 32)
	if _, err = io.ReadFull(rand.Reader, b); err != nil {
		return
	}
	id = base64.URLEncoding.EncodeToString(b)
	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	_, err := r.Cookie("session_id")
	if err == nil {
		// user has already logged in
		return
	}

	u := &User{}
	err = cleanUser(r, u)
	if invalid := u.User == "" || u.Password == ""; err != nil || invalid {
		if invalid || err.Error() == "json error" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Println(err, "in loginHandler in getUserFromRequestBody")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: go to database, ask for user...
	// db name????? select???? request??? extra api?
	// after db request:

	// stub:
	dbResponse := User{
		"test",
		"test",
	}

	if u.User == dbResponse.User && u.Password == dbResponse.Password { // stub
		sessionID, err := generateSessionID()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		cookie := http.Cookie{
			Name:  "session_id",
			Value: sessionID,
			// Expires: time.Now().Add(24 * time.Hour),
		}
		http.SetCookie(w, &cookie)
		log.Println("User logged in:", u.User)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		// user has already logged out
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	args := os.Args[1:]
	if len(args) != 1 {
		panic("no port is given or more than two args are given")
	}
	port := args[0]

	if portNum, err := strconv.Atoi(port); err != nil {
		panic("wrong port format is given")
	} else if portNum < 0 || portNum > 65535 {
		panic("port should be ≥ 0 and ≤ 65535")
	}

	log.Println("starting server at:", port)
	http.ListenAndServe(":"+port, nil)
}
