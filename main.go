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

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", middleware(Index))
	http.HandleFunc("/signup", middleware(Signup))
	http.HandleFunc("/login", middleware(Login))
	http.HandleFunc("/logoff", Logoff)
	http.HandleFunc("/account", auth(Account))
	http.HandleFunc("/post", auth(Post))

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
