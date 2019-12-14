package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/imtanmoy/authy/config"
	"github.com/imtanmoy/authy/db"
	"github.com/imtanmoy/authy/logger"
)

// Server provides an http.Server.
type Server struct {
	*http.Server
}

// NewServer creates and configures an APIServer serving all application routes.
func NewServer() (*Server, error) {
	logger.Info("configuring server...")
	handler, err := New()
	if err != nil {
		return nil, err
	}

	host := config.Conf.SERVER.HOST
	port := strconv.Itoa(config.Conf.SERVER.PORT)
	addr := host + ":" + port

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	srv := http.Server{
		Addr:         addr,
		Handler:      (middlewares{tracing(nextRequestID), logging()}).apply(handler),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{&srv}, nil
}

// Start runs ListenAndServe on the http.Server with graceful shutdown.
func (srv *Server) Start(ctx context.Context) (err error) {
	logger.Info("starting server...")
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatalf("listen:%+s\n", err)
		}
	}()
	logger.Infof("Listening on %s\n", srv.Addr)

	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	srv.SetKeepAlivesEnabled(false)
	if err = srv.Shutdown(ctxShutDown); err != nil {
		logger.Fatalf("server Shutdown Failed:%+s", err)
	}
	logger.Info("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}
	dbErr := db.Shutdown()
	if dbErr != nil {
		logger.Errorf("%s : %s", "Database shutdown failed", dbErr)
	}
	return
}
