package main

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var flashStore = sessions.NewCookieStore([]byte("a-secret-string"))
var flashMessageKey = "flash-message"
var flashMessageName = "flash-session"

func SetFlash(w http.ResponseWriter, r *http.Request, msg string) {
	session, err := flashStore.Get(r, flashMessageName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.AddFlash(msg, flashMessageKey)
	session.Save(r, w)
}

func GetFlash(w http.ResponseWriter, r *http.Request) ([]any, error) {

	session, err := flashStore.Get(r, flashMessageName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	if flashes := session.Flashes(flashMessageKey); len(flashes) > 0 {
		session.Save(r, w)
		return flashes, nil
	}
	return nil, err
}
