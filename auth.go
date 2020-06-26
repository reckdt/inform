package main

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// auth
func auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, sessionId, err := getUsernameAndSessionId(r)
		if err != nil || !verifySession(sessionId, username) {
			http.Redirect(w, r, "/login", 302)
		} else {
			h.ServeHTTP(w, r)
		}
	}
}

func addCookie(w http.ResponseWriter, name string, value string, ttl time.Duration) {
	expires := time.Now().Add(ttl)
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expires,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func removeCookie(w http.ResponseWriter, name string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func getHashedPassword(password string) string {
	bytePassword := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

func compareHashAndPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	} else {
		return true
	}
}

func getSessionId() string {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(out))
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
