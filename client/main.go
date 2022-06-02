package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func main() {

	fmt.Println("Client started at 8080")

	mux := http.NewServeMux()
	tmpl = template.Must(template.ParseGlob("templates/*.gohtml"))
	fs := http.FileServer(http.Dir("./static"))

	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", login)
	mux.HandleFunc("/enterProfileInformation", enterProfileInformation)
	mux.HandleFunc("/signUp", signUp)
	mux.HandleFunc("/mainProfile", mainProfile)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
