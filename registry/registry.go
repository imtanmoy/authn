package registry

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/imtanmoy/authn/config"
	"github.com/imtanmoy/authn/events"
	"github.com/imtanmoy/logx"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

type Registry interface {
	Init() error
	Config() config.Config
	Bus() events.EventBus
	DB() *sql.DB
	Close()
}

var _ Registry = (*registry)(nil)

type registry struct {
	c  config.Config
	b  events.EventBus
	db *sql.DB
}

func (r *registry) Config() config.Config {
	return r.c
}

func (r *registry) Bus() events.EventBus {
	if r.b == nil {
		r.b = events.New()
	}
	return r.b
}

func (r *registry) DB() *sql.DB {
	if r.db == nil {
		db, err := connectDB(r.c.DB.HOST, r.c.DB.PORT, r.c.DB.USERNAME, r.c.DB.PASSWORD, r.c.DB.DBNAME)
		if err != nil {
			logx.Fatalf("%s : %s", "Database Could not be initiated", err)
		}
		logx.Info("Database Initiated...")
		r.db = db
	}
	return r.db
}

func NewRegistry(c config.Config) Registry {
	return &registry{c: c}
}

func (r *registry) Init() error {
	bus := events.New()
	r.b = bus
	db, err := connectDB(r.c.DB.HOST, r.c.DB.PORT, r.c.DB.USERNAME, r.c.DB.PASSWORD, r.c.DB.DBNAME)
	if err != nil {
		return err
	}
	r.db = db
	return nil
}

func (r *registry) Close() {
	err := r.db.Close()
	if err != nil {
		logx.Errorf("%s : %s", "Database shutdown failed", err)
	}
}

func connectDB(host string, port int, username, password, database string) (*sql.DB, error) {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", username, password, host, port, database)
	db := ConnectDBViaPgx(connString)
	return db, nil
}

func ConnectPgx(connString string) (*pgx.Conn, error) {
	connConfig, err := pgx.ParseConfig(connString)
	if err != nil {
		panic(err)
	}
	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	return conn, err
}

func ConnectDBViaPgx(connString string) *sql.DB {
	connConfig, err := pgx.ParseConfig(connString)
	if err != nil {
		panic(err)
	}
	return stdlib.OpenDB(*connConfig)
}
