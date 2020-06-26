package main

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
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
