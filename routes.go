package main

import (
	"fmt"
	"net/http"
)

const title string = "inform.lol"

// index
func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	p := AuthPage{
		Title: title,
	}

	templates.ExecuteTemplate(w, "index.html", p)
}

// error handler
func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404 Page Not Found.")
	}
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
		username, password := getUserNameAndPassword(r)
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

		createSessionAndCookies(w, username)
		http.Redirect(w, r, "/", 302)
	}
}

// login
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		p := AuthPage{Title: title}
		templates.ExecuteTemplate(w, "login.html", p)
	} else if r.Method == "POST" {
		username, password := getUserNameAndPassword(r)
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

		createSessionAndCookies(w, username)
		http.Redirect(w, r, "/", 302)
	}
}

func getUserNameAndPassword(r *http.Request) (string, string) {
	r.ParseForm()
	return r.PostForm.Get("username"), r.PostForm.Get("password")
}

//logoff
func Logoff(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		removeSessionAndCookies(w, r)
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
		username := r.Context().Value("username").(string)
		p := AccountPage{
			Title:    title,
			Username: username,
		}
		templates.ExecuteTemplate(w, "account.html", p)
	}
}
