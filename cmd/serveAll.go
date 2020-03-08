package cmd

import (
	"github.com/imtanmoy/authn/registry"
	"github.com/imtanmoy/authn/server/http"
	"github.com/imtanmoy/logx"
	"sync"
)

func serveHttp(wg *sync.WaitGroup, r registry.Registry) {
	defer wg.Done()
	srv, err := http.NewServer(r)
	if err != nil {
		logx.Fatalf("%s : %s", "Server could not be started", err)
	}
	srv.Run()
}

func startEventBus(wg *sync.WaitGroup, r registry.Registry) {
	defer wg.Done()
	r.Bus().Run()
}

func ServeAll(r registry.Registry) {
	var wg sync.WaitGroup
	wg.Add(2)
	go serveHttp(&wg, r)
	go startEventBus(&wg, r)
	wg.Wait()
}
