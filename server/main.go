package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	fmt.Println("Server started at 9090")

	mux := http.NewServeMux()

	mux.HandleFunc("/authLogin", AuthLogin)
	mux.HandleFunc("/addNewUser", AddNewUser)
	mux.HandleFunc("/enterProfileInformation", EnterProfileInformation)
	mux.HandleFunc("/mainProfile", MainProfile)
	mux.HandleFunc("/logout", Logout)
	mux.HandleFunc("/googleSession", GoogleSession)
	mux.HandleFunc("/callback", Callback)

	log.Fatal(http.ListenAndServe(":9090", mux))
}
