package main

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"os/exec"
	"strings"
	"time"
)

func addUser(username string, password string) error {
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

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errors.New("Incorrect password.")
	}

	return nil
}

func getHashedPassword(password string) string {
	bytePassword := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

func getSessionId() string {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(out))
}

func addSession(sessionId string, username string) error {
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
