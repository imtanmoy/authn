package tests

import (
	"database/sql"
	"fmt"
	"github.com/imtanmoy/authn/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
	"strings"
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

func FakeUsers(nums int) []*models.User {
	var users []*models.User
	for i := 0; i < nums; i++ {
		str := strconv.Itoa(i)
		u := &models.User{
			Name:     "Test " + str,
			Email:    "test" + str + "@test.com",
			Password: "password",
		}
		users = append(users, u)
	}
	return users
}

func SeedUsers(db *sql.DB) {
	users := FakeUsers(10)
	err := InsertTestUsers(db, users)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertTestUsers(db *sql.DB, users []*models.User) error {
	valueStrings := make([]string, 0, len(users))
	valueArgs := make([]interface{}, 0, len(users)*3)
	for i, u := range users {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))

		valueArgs = append(valueArgs, u.Name)
		valueArgs = append(valueArgs, u.Email)
		valueArgs = append(valueArgs, u.Password)
	}
	smt := `INSERT INTO users(name, email, password) VALUES %s`
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(smt, valueArgs...)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func FakeOrgs(nums int) []*models.Organization {
	var orgs []*models.Organization
	for i := 0; i < nums; i++ {
		str := strconv.Itoa(i)
		u := &models.Organization{
			Name:    "Test orgs " + str,
			OwnerID: 1,
		}
		orgs = append(orgs, u)
	}
	return orgs
}

func InsertTestOrgs(db *sql.DB, orgs []*models.Organization) error {
	valueStrings := make([]string, 0, len(orgs))
	valueArgs := make([]interface{}, 0, len(orgs)*2)
	for i, org := range orgs {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))

		valueArgs = append(valueArgs, org.Name)
		valueArgs = append(valueArgs, org.OwnerID)
	}
	smt := `INSERT INTO organizations(name, owner_id) VALUES %s`
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(smt, valueArgs...)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
