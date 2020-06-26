package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"time"
)

var db *sql.DB
var templates *template.Template

func main() {
	templates = template.Must(template.ParseGlob("templates/*.html"))

	http.HandleFunc("/", Index)
	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/account", Account)

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

func addCookie(w http.ResponseWriter, name string, value string, ttl time.Duration) {
	expires := time.Now().Add(ttl)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expires,
	}
	http.SetCookie(w, &cookie)
}
