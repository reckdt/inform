package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
)

var db *sql.DB
var templates *template.Template

func main() {
	templates = template.Must(template.ParseGlob("templates/*.html"))

	http.HandleFunc("/", Index)
	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/logoff", Logoff)
	http.HandleFunc("/account", auth(Account))

	initDb()

	fmt.Println("Server running...")

	log.Fatal(http.ListenAndServe(":9990", nil))
}

func initDb() {
	var err error
	db, err = sql.Open("postgres", "user=postgres password=postgres dbname=mydb sslmode=disable")
	if err != nil {
		panic(err)
	}
}
