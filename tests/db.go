package tests

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/imtanmoy/authn/registry"
	"github.com/imtanmoy/logx"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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

func ConnectionPool(connString string) *pgxpool.Pool {
	poolConfig, err := pgxpool.ParseConfig(connString + "?pool_max_conns=10")
	if err != nil {
		logx.Fatal(err)
	}
	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		logx.Println("acquired a connection from the pool")
		return conn.Ping(ctx)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		logx.Fatal(err)
	}
	return pool
}
