package tests

import (
	"github.com/go-pg/pg/v9"
	"log"
)


func ConnectTestDB(username, password, database, address string) (*pg.DB, error) {
	connect := pg.Connect(&pg.Options{
		User:     username,
		Password: password,
		Database: database,
		Addr:     address,
	})
	return connect, nil
}

func TruncateTestDB(db *pg.DB)  {
	_, err := db.Exec("TRUNCATE TABLE users, organizations, invitations, users_organizations RESTART IDENTITY;")
	if err != nil {
		log.Fatal(err)
	}
}
