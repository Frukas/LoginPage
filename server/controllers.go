package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"loginpage/DBhandler"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

//Todo Make as separated package - separete in diferent files to improve organization

//Todo generate random number
var randomState = "Random"
var db DBhandler.DB

var endconnectionString = "http://" + endconnection + ":9090"
var frontconnectionString = "http://" + frontconnection + ":8080"

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  endconnectionString + "/callback",
	ClientID:     "XXXXXXXXXXXXXXXXXXXXXXXXXX.apps.googleusercontent.com",
	ClientSecret: "XXXXXX-XXX-XXXX-XXXXXXXXXXXXXXXX",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

func AuthLogin(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		res := db.GetEmailDetails(r.FormValue("email"))

		if res.Email == "" {
			http.Redirect(w, r, frontconnectionString+"/login?message=Username entered does not exist", http.StatusSeeOther)
		}

		if res.Password == r.FormValue("password") {

			startSession(w, r)

			http.Redirect(w, r, frontconnectionString+"/mainProfile", http.StatusSeeOther)
		} else {

			http.Redirect(w, r, frontconnectionString+"/login?message=Password is incorrect", http.StatusSeeOther)
		}
	}
}

func AddNewUser(w http.ResponseWriter, r *http.Request) {
	isgoogle := false

	session, err := store.Get(r, "mysession")
	if err != nil {
		log.Println("Error getting the session: ", err)
	}

	if session.Values["googleMail"] != nil {
		isgoogle = true
	}

	res := db.GetEmailDetails(r.FormValue("email"))

	redirectPage := frontconnectionString + "/enterProfileInformation"

	user := DBhandler.Users{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password1"),
		IsGoogle: isgoogle,
	}
	if res.Email == user.Email {
		http.Redirect(w, r, frontconnectionString+"?message=This E-Mail is already registed", http.StatusSeeOther)
	} else {
		db.AddNewUser(user)

		startSession(w, r)

		http.Redirect(w, r, redirectPage, http.StatusSeeOther)
	}
}

func EnterProfileInformation(w http.ResponseWriter, r *http.Request) {

	user := DBhandler.Users{
		Name:      r.FormValue("username"),
		Email:     r.FormValue("email"),
		Telephone: r.FormValue("telephone"),
	}
	session, err := store.Get(r, "mysession")
	if err != nil {
		log.Println("Error getting the session: ", err)
	}

	userN := strings.TrimSuffix(fmt.Sprintf("%v", session.Values["username"]), "/")

	db.UpdateUser(user, userN)

	startSession(w, r)

	redirectPage := frontconnectionString + "/mainProfile"

	http.Redirect(w, r, redirectPage, http.StatusSeeOther)
}

func MainProfile(w http.ResponseWriter, r *http.Request) {

	resp := r.URL.Query()["userEmail"]

	email := stripSuffixPrefix(resp[0])

	details := db.GetEmailDetails(email)

	json.NewEncoder(w).Encode(details)
}

func Logout(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "mysession")
	if err != nil {
		log.Println("Error getting the session: ", err)
	}
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redirectPage := frontconnectionString + "/"
	http.Redirect(w, r, redirectPage, http.StatusSeeOther)
}

func GoogleSession(w http.ResponseWriter, r *http.Request) {

	url := googleOauthConfig.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

//todo break into smaller functions - Separation of concerns
func Callback(w http.ResponseWriter, r *http.Request) {
	var gaccount googleAccountInfo
	if r.FormValue("state") != randomState {
		log.Println("No mach state")
		return
	}
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		log.Println("Error retriving token: ", err)
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Println("Error getting data from google: ", err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading the body: ", err)
		return
	}
	json.Unmarshal(body, &gaccount)

	dbcheck := db.GetEmailDetails(gaccount.Email)

	if dbcheck.Email != "" {
		session, err := store.Get(r, "mysession")
		if err != nil {
			log.Println("Error getting the session: ", err)
		}
		session.Values["username"] = gaccount.Email

		session.Options.MaxAge = 180
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, frontconnectionString+"/enterProfileInformation", http.StatusTemporaryRedirect)

	} else {
		session, err := store.Get(r, "mysession")
		if err != nil {
			log.Println("Error getting the session: ", err)
		}

		session.Values["googleMail"] = gaccount.Email

		session.Options.MaxAge = 180
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, frontconnectionString+"/signUp", http.StatusTemporaryRedirect)
	}
}
