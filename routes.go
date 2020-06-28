package main

import (
	"fmt"
	"net/http"
)

type Page struct {
	Username string
}

func (p Page) Title() string {
	return "inform.lol"
}

// index
func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	p := AuthPage{}
	p.Username = getUsername(r)

	posts := getPosts()
	for _, p := range posts {
		println(p.Url)
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
	Page
	ErrMsg   string
	Username string
	Password string
}

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		p := AuthPage{}
		templates.ExecuteTemplate(w, "signup.html", p)
	} else if r.Method == "POST" {
		username, password := getUserNameAndPassword(r)
		err := addUser(username, password)
		if err != nil {
			p := AuthPage{
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
		p := AuthPage{}
		templates.ExecuteTemplate(w, "login.html", p)
	} else if r.Method == "POST" {
		username, password := getUserNameAndPassword(r)
		err := verifyUser(username, password)
		if err != nil {
			p := AuthPage{
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
	Page
}

func Account(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		p := AccountPage{}
		p.Username = getUsername(r)
		templates.ExecuteTemplate(w, "account.html", p)
	}
}

// post
type PostPage struct {
	Page
	ErrMsg    string
	PostTitle string
	Url       string
	Text      string
}

func Post(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		p := PostPage{}
		p.Username = getUsername(r)
		templates.ExecuteTemplate(w, "post.html", p)
	} else if r.Method == "POST" {
		r.ParseForm()
		title := r.PostForm.Get("title")
		url := r.PostForm.Get("url")
		text := r.PostForm.Get("text")
		post := UserPost{
			Title:    trim(title),
			Url:      trim(url),
			Text:     trim(text),
			Username: getUsername(r),
		}
		err := addPost(post)
		if err != nil {
			p := PostPage{
				ErrMsg:    err.Error(),
				PostTitle: title,
				Url:       url,
				Text:      text,
			}
			templates.ExecuteTemplate(w, "post.html", p)
			return
		}
		http.Redirect(w, r, "/", 302)
	}
}

// gets username from context
func getUsername(r *http.Request) string {
	username, ok := r.Context().Value("username").(string)
	if !ok {
		return ""
	}
	return username
}
