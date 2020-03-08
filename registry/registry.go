package registry

import (
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authn/config"
	"github.com/imtanmoy/authn/events"
	"github.com/imtanmoy/logx"
	"strconv"
)

type Registry interface {
	Init() error
	Config() config.Config
	Bus() events.EventBus
	DB() *pg.DB
	Close()
}

var _ Registry = (*registry)(nil)

type registry struct {
	c  config.Config
	b  events.EventBus
	db *pg.DB
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

func (r *registry) DB() *pg.DB {
	if r.db == nil {
		db, err := connectDB(r.c.DB.USERNAME, r.c.DB.PASSWORD, r.c.DB.DBNAME, r.c.DB.HOST+":"+strconv.Itoa(r.c.DB.PORT))
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
	db, err := connectDB(r.c.DB.USERNAME, r.c.DB.PASSWORD, r.c.DB.DBNAME, r.c.DB.HOST+":"+strconv.Itoa(r.c.DB.PORT))
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

func connectDB(username, password, database, address string) (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		User:     username,
		Password: password,
		Database: database,
		Addr:     address,
	})
	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		return nil, err
	}
	return db, nil
}
