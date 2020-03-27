package repository

import (
	"database/sql"
	"github.com/imtanmoy/authn/organization"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"log"
)

type repository struct {
	conn *pgx.Conn
}

var _ organization.Repository = (*repository)(nil)

// NewRepository will create an object that represent the organization.Repository interface
func NewRepository(db *sql.DB) organization.Repository {
	conn, err := stdlib.AcquireConn(db)
	if err != nil {
		log.Fatal(err)
	}
	return &repository{conn: conn}
}
