package http

import (
	"github.com/go-chi/chi"
	_chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/imtanmoy/authn/config"
	"github.com/imtanmoy/authn/registry"
	"github.com/imtanmoy/logx"
	"github.com/ory/graceful"
	"net/http"
	"strconv"
	"time"
)

// Server provides an http.Server.
type Server struct {
	*http.Server
}

func newRouter() (*chi.Mux, error) {

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

// NewServer creates and configures an APIServer serving all application routes.
func NewServer(r registry.Registry) (*Server, error) {
	logx.Info("configuring server...")
	handler, err := newRouter()
	if err != nil {
		return nil, err
	}
	RegisterHandler(handler, r)

	host := config.Conf.SERVER.HOST
	port := strconv.Itoa(config.Conf.SERVER.PORT)
	addr := host + ":" + port

	server := graceful.WithDefaults(&http.Server{
		Addr:    addr,
		Handler: handler,
	})

	return &Server{server}, nil
}

// Run runs ListenAndServe on the http.Server with graceful shutdown.
func (server *Server) Run() {
	logx.Printf("Starting the httpd on: %s", server.Addr)
	if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
		logx.Fatal("Failed to gracefully shutdown httpd")
	}
	return
}
