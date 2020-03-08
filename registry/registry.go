package registry

import (
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authn/config"
	"github.com/imtanmoy/authn/db"
	"github.com/imtanmoy/authn/events"
)

type Registry interface {
	Init() error
	Config() config.Config
	Bus() events.EventBus
	DB() *pg.DB
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
		r.db = db.DB
	}
	return r.db
}

func NewRegistry(c config.Config) Registry {
	return &registry{c: c}
}

func (r *registry) Init() error {
	bus := events.New()
	r.b = bus
	return nil
}
