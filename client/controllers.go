package main

//Todo Make as separated package - separate according to the templates

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
)

var tmpl *template.Template

var store = sessions.NewCookieStore([]byte("mysession"))

func login(w http.ResponseWriter, r *http.Request) {

	data := passData{
		Message: r.URL.Query().Get("message"),
	}

	if err := tmpl.ExecuteTemplate(w, "login.gohtml", data); err != nil {
		fmt.Fprintf(w, "Error loading login page: %v", err)
	}
}

func enterProfileInformation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	session, err := store.Get(r, "mysession")
	if err != nil {
		log.Println("Error getting the session: ", err)
	}

	if session.Values["username"] == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	userDetails := getUserInfo(convertInterfaceToString(session.Values["username"]))

	data := passData{
		Message: "",
		User:    userDetails,
	}

	if err := tmpl.ExecuteTemplate(w, "enterProfileInformation.gohtml", data); err != nil {
		fmt.Fprintf(w, "Error loading enterProfileInformation page: %v", err)
	}
}

func signUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	email := ""

	session, err := store.Get(r, "mysession")
	if err != nil {
		log.Println("Error getting the session: ", err)
	}

	if convertInterfaceToString(session.Values["googleMail"]) != "<nil>" {
		email = convertInterfaceToString(session.Values["googleMail"])
	}

	msg := r.URL.Query().Get("message")
	data := passData{
		Message: msg,
		User:    user{Email: email},
	}

	if err := tmpl.ExecuteTemplate(w, "singup.gohtml", data); err != nil {
		fmt.Fprintf(w, "Error loading singup page: %v", err)
	}
}

func mainProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	session, err := store.Get(r, "mysession")
	if err != nil {
		log.Println("Error getting the session: ", err)
	}

	if session.Values["username"] == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	userDetails := getUserInfo(convertInterfaceToString(session.Values["username"]))

	if err := tmpl.ExecuteTemplate(w, "mainProfile.gohtml", userDetails); err != nil {
		fmt.Fprintf(w, "Error loading mainProfil page: %v", err)
	}
}
