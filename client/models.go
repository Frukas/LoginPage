package main

//Todo Make as separated package

type passData struct {
	Message  string	
	User     user
}

type user struct {
	Name      string
	Email     string
	Telephone string
	IsGoogle  bool
}