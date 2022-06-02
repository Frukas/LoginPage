package main

//Todo Make as separated packagego run 

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func convertInterfaceToString(inter interface{}) string {
	return fmt.Sprintf("%v", inter)
}

func getUserInfo(userEmail string) user {

	var us user
	getRequestString := "http://localhost:9090/mainProfile?userEmail=" + userEmail
	resp, err := http.Get(getRequestString)
	if err != nil {
		log.Fatal("Error getting the user info", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &us)

	return us
}
