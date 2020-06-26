package main

import (
	"errors"
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

func getUsernameAndSessionId(r *http.Request) (string, string, error) {
	cookie, err := r.Cookie("username")
	if err != nil {
		return "", "", errors.New("No username cookie")
	}
	username := cookie.Value

	cookie, err = r.Cookie("sessionId")
	if err != nil {
		return "", "", errors.New("No session id cookie")
	}
	sessionId := cookie.Value

	return username, sessionId, nil
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

func createSessionAndCookies(w http.ResponseWriter, username string) {
	sessionId := getSessionId()
	addSession(sessionId, username)

	addCookie(w, "username", username, 60*time.Minute)
	addCookie(w, "sessionId", sessionId, 60*time.Minute)
}

func removeSessionAndCookies(w http.ResponseWriter, r *http.Request) {
	username, sessionId, err := getUsernameAndSessionId(r)
	if err == nil {
		removeSession(sessionId, username)
	}
	removeCookie(w, "username")
	removeCookie(w, "sessionId")
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
		p := AccountPage{
			Title:    title,
			Username: "ryan",
		}
		templates.ExecuteTemplate(w, "account.html", p)
	}
}
