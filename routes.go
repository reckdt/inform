package main

import (
	"fmt"
	"net/http"
	"time"
)

const title string = "Inform.lol"

// index
func Index(w http.ResponseWriter, r *http.Request) {

	//	cookie, err := r.Cookie("username")
	//	if err != nil {
	//		fmt.Println("no cookie")
	//	} else {
	//		fmt.Println(cookie.Value)
	//	}

	fmt.Fprintf(w, "This is the index.")
}

// signup
type AuthPage struct {
	Title    string
	ErrMsg   string
	Username string
	Password string
}

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		p := AuthPage{Title: title}
		templates.ExecuteTemplate(w, "signup.html", p)
	} else if r.Method == "POST" {
		r.ParseForm()
		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")

		err := addUser(username, password)
		if err != nil {
			p := AuthPage{
				Title:    title,
				ErrMsg:   err.Error(),
				Username: username,
				Password: password,
			}
			templates.ExecuteTemplate(w, "signup.html", p)
			return
		}

		addCookie(w, "username", username, 30*time.Minute)
		addCookie(w, "password", password, 30*time.Minute)

		fmt.Println(username, password)
		http.Redirect(w, r, "/", 302)
	}
}

// login
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		p := AuthPage{Title: title}
		templates.ExecuteTemplate(w, "login.html", p)
	} else if r.Method == "POST" {
		r.ParseForm()
		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")

		err := verifyUser(username, password)
		if err != nil {
			p := AuthPage{
				Title:    title,
				ErrMsg:   err.Error(),
				Username: username,
				Password: password,
			}
			templates.ExecuteTemplate(w, "login.html", p)
			return
		}

		addCookie(w, "username", username, 30*time.Minute)
		addCookie(w, "password", password, 30*time.Minute)

		fmt.Println("logged in")
		http.Redirect(w, r, "/", 302)
	}
}

// account

type AccountPage struct {
	Title    string
	Username string
}

func Account(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		p := AccountPage{
			Title:    title,
			Username: "ryan",
		}
		templates.ExecuteTemplate(w, "account.html", p)
	}
}
