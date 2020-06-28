package main

import (
	"errors"
	"net/url"
	"strings"
	"time"
)

type UserPost struct {
	Id       int
	Username string
	Title    string
	Url      string
	Text     string
	Posted   time.Time
	Host     string
}

type UserComment struct {
	Id        int
	PostId    int
	CommentId int
	Text      string
	Posted    time.Time
}

func empty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func trim(s string) string {
	return strings.TrimSpace(s)
}

func addPost(p UserPost) error {
	p.Title = trim(p.Title)
	p.Url = trim(p.Url)
	p.Text = trim(p.Text)

	if len(p.Title) < 5 {
		return errors.New("Title must be atleast 5 characters.")
	}

	if len(p.Title) > 100 {
		return errors.New("Title must be less than 100 characters.")
	}

	if empty(p.Text) && empty(p.Url) {
		return errors.New("Url or Text must be entered.")
	}

	if len(p.Text) > 10000 {
		return errors.New("Text must be less than 10,000 characters.")
	}

	_, err := url.ParseRequestURI(trim(p.Url))
	if err != nil {
		return errors.New("Invalid Url.")
	}

	_, err = db.Exec("INSERT INTO posts (username, title, url, text, dt) VALUES ($1, $2, $3, $4, $5)",
		p.Username, p.Title, p.Url, p.Text, time.Now())
	if err != nil {
		return errors.New("Internal error making post.")
	}

	return nil
}

func getPosts() []UserPost {
	posts := []UserPost{}
	rows, err := db.Query("SELECT * from posts")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var p UserPost
		err := rows.Scan(&p.Id, &p.Username, &p.Title, &p.Url, &p.Text, &p.Posted)
		if err != nil {
			panic(err)
		} else {
			posts = append(posts, p)
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return posts
}

func addUser(username string, password string) error {
	username = trim(username)

	if len(username) < 3 {
		return errors.New("Username must be atleast 3 characters.")
	}

	if len(password) < 8 {
		return errors.New("Password must be atleast 8 characters.")
	}

	_, err := db.Exec("INSERT INTO users VALUES ($1, $2)", username, getHashedPassword(password))
	if err != nil {
		return errors.New("Username is already taken.")
	}

	return nil
}

func verifyUser(username string, password string) error {
	var hashedPassword string
	row := db.QueryRow("SELECT username, password FROM users WHERE username = $1", username)
	err := row.Scan(&username, &hashedPassword)

	if err != nil {
		return errors.New("Incorrect username.")
	}

	if !compareHashAndPassword(hashedPassword, password) {
		return errors.New("Incorrect password.")
	}

	return nil
}

func verifySession(sessionId string, username string) bool {
	row := db.QueryRow("SELECT id, username FROM sessions WHERE id = $1 AND username = $2", sessionId, username)
	err := row.Scan(&sessionId, &username)

	if err != nil {
		return false
	} else {
		return true
	}
}

func addSession(sessionId string, username string) error {
	removeAllUserSessions(username)
	_, err := db.Exec("INSERT INTO sessions VALUES ($1, $2, $3)", sessionId, username, time.Now())
	if err != nil {
		return errors.New("Error adding session.")
	}
	return nil
}

func removeSession(sessionId string, username string) error {
	_, err := db.Exec("DELETE FROM sessions WHERE id = $1 AND username = $2", sessionId, username)
	if err != nil {
		return errors.New("Error removing session.")
	}
	return nil
}

func removeAllUserSessions(username string) error {
	_, err := db.Exec("DELETE FROM sessions WHERE username = $1", username)
	if err != nil {
		return errors.New("Error removing session.")
	}
	return nil
}
