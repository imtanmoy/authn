package server

import (
	"github.com/go-chi/chi"
	_chiMiddleware "github.com/go-chi/chi/middleware"
	"time"
)

func NewRouter() (*chi.Mux, error) {

	r := chi.NewRouter()
	r.Use(_chiMiddleware.Recoverer)
	r.Use(_chiMiddleware.RequestID)
	r.Use(_chiMiddleware.RealIP)
	r.Use(_chiMiddleware.DefaultCompress)
	r.Use(_chiMiddleware.Timeout(15 * time.Second))
	r.Use(_chiMiddleware.Logger)
	r.Use(_chiMiddleware.AllowContentType("application/json"))
	r.Use(_chiMiddleware.Heartbeat("/heartbeat"))
	//r.Use(render.SetContentType(3)) //render.ContentTypeJSON resolve value 3
	return r, nil
}
