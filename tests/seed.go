package tests

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func SeedUser(db *sql.DB) {
	hashedPass, _ := hashPassword("password")
	name := "Test"
	email := "test@test.com"
	password := hashedPass
	_, err := db.Exec("INSERT INTO users(name, email, password) VALUES ($1,$2,$3)", name, email, password)
	if err != nil {
		log.Fatal(err)
	}
}

func SeedOrganization(db *sql.DB) {
	name := "Test Organization"
	_, err := db.Exec("INSERT INTO organizations(name) VALUES ($1)", name)
	if err != nil {
		log.Fatal(err)
	}
}
