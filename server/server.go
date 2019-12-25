package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/imtanmoy/authN/config"
	"github.com/imtanmoy/authN/db"
	"github.com/imtanmoy/logx"
)

// Server provides an http.Server.
type Server struct {
	*http.Server
}

// NewServer creates and configures an APIServer serving all application routes.
func NewServer() (*Server, error) {
	logx.Info("configuring server...")
	handler, err := New()
	if err != nil {
		return nil, err
	}

	host := config.Conf.SERVER.HOST
	port := strconv.Itoa(config.Conf.SERVER.PORT)
	addr := host + ":" + port

	srv := http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{&srv}, nil
}

// Run runs ListenAndServe on the http.Server with graceful shutdown.
func (srv *Server) Run(ctx context.Context) (err error) {
	logx.Info("starting server...")
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logx.Fatalf("listen:%+s\n", err)
		}
	}()
	logx.Infof("Listening on %s\n", srv.Addr)

	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	srv.SetKeepAlivesEnabled(false)
	if err = srv.Shutdown(ctxShutDown); err != nil {
		logx.Fatalf("server Shutdown Failed:%+s", err)
	}
	logx.Info("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}
	dbErr := db.Shutdown()
	if dbErr != nil {
		logx.Errorf("%s : %s", "Database shutdown failed", dbErr)
	}
	return
}
