package cmd

import (
	"context"
	"github.com/imtanmoy/authn/registry"
	"github.com/imtanmoy/authn/server"
	"github.com/imtanmoy/logx"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func serveHttp(wg *sync.WaitGroup, r registry.Registry) {
	defer wg.Done()
	// initializing newServer
	newServer, err := server.NewServer(r)
	if err != nil {
		logx.Fatalf("%s : %s", "Server could not be started", err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGSTOP)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		logx.Infof("system call:%+v", oscall)
		cancel()
	}()

	if err := newServer.Run(ctx); err != nil {
		logx.Infof("failed to serve:+%v\n", err)
	}
	close(c)
}

func startEventBus(wg *sync.WaitGroup, r registry.Registry) {
	defer wg.Done()
	r.Bus().Run(context.Background())
}

func ServeAll(r registry.Registry) {
	var wg sync.WaitGroup
	wg.Add(2)
	go serveHttp(&wg, r)
	go startEventBus(&wg, r)
	wg.Wait()
}
