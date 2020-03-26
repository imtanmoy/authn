package tests

import (
	"database/sql"
	"fmt"
	"github.com/imtanmoy/authn/registry"
	"log"
)

func ConnectTestDB(host string, port int, username, password, database string) (*sql.DB, error) {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", username, password, host, port, database)
	db := registry.ConnectDBViaPgx(connString)
	return db, nil
}

func TruncateTestDB(db *sql.DB) {
	_, err := db.Exec("TRUNCATE TABLE users, organizations, invitations, users_organizations RESTART IDENTITY;")
	if err != nil {
		log.Fatal(err)
	}
}
