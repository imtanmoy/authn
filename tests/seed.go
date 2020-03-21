package tests

import (
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authn/models"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func SeedUser(db *pg.DB) {
	hashedPass, _ := hashPassword("password")
	var u models.User
	u.Name = "Test"
	u.Email = "tests@tests.com"
	u.Password = hashedPass
	err := db.Insert(u)
	if err != nil {
		log.Fatal(err)
	}
}

func SeedOrganization(db *pg.DB) {
	var org models.Organization
	org.Name = "Test Organization"
	err := db.Insert(org)
	if err != nil {
		log.Fatal(err)
	}
}
