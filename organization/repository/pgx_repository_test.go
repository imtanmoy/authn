package repository

import (
	"database/sql"
	"github.com/imtanmoy/authn/organization"
	"github.com/imtanmoy/authn/tests"
	"log"
)

var db *sql.DB
var repo organization.Repository

func init() {
	var err error
	db, err = tests.ConnectTestDB("localhost", 5432, "admin", "password", "authn")
	if err != nil {
		log.Fatal(err)
	}
	repo = NewRepository(db)
}
