package DBhandler

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Users struct {
	Id        int    `json:"Id"`
	Name      string `json:"Name"`
	Email     string `json:"Email"`
	Password  string `json:"Password"`
	Telephone string `json:"Telephone"`
	IsGoogle  bool   `json:"IsGoogle"`
}

type DB struct {
	db *sql.DB
}

func (d *DB) init() {
	d.db = getDBconnection()
}

func getDBconnection() *sql.DB {
	db, err := sql.Open("mysql", "b5cb5c82f5f8e9:afd8a4c0@tcp(us-cdbr-east-05.cleardb.net:3306)/heroku_ae713263b61ed52")

	if err != nil {
		panic(err.Error())
	}

	return db
}

func (d *DB) GetEmailDetails(name string) Users {
	d.init()
	defer d.Close()
	var user Users

	res, err := d.db.Query("SELECT * FROM heroku_ae713263b61ed52.loginpage where email = ?", name)
	if err != nil {
		panic(err.Error())
	}
	res.Next()
	res.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Telephone, &user.IsGoogle)

	return user
}

func (d *DB) AddNewUser(user Users) {
	d.init()
	defer d.Close()
	res, err := d.db.Prepare("INSERT Into heroku_ae713263b61ed52.loginpage (name, email, password, telephone, isgoogle) Values(?,?,?,?, ?)")
	if err != nil {
		panic(err.Error())
	}
	_, err = res.Exec(user.Name, user.Email, user.Password, user.Telephone, user.IsGoogle)
	if err != nil {
		panic(err.Error())
	}
}

func (d *DB) UpdateUser(user Users, oldmail string) {
	d.init()
	defer d.Close()

	fmt.Println(oldmail)

	res, err := d.db.Prepare("UPDATE heroku_ae713263b61ed52.loginpage SET name = ?, email = ?, telephone = ? where email = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = res.Exec(user.Name, user.Email, user.Telephone, oldmail)
	if err != nil {
		panic(err.Error())
	}
}

func (d *DB) Close() {
	d.db.Close()
}
