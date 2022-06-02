package main

import (
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

//Todo Make as separated package

var store = sessions.NewCookieStore([]byte("mysession"))

func stripSuffixPrefix(value string) string {
	value = strings.TrimPrefix(value, "[")
	value = strings.TrimSuffix(value, "]")
	return value
}

func startSession(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mysession")
	session.Values["username"] = r.FormValue("email")

	session.Options.MaxAge = 180
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
