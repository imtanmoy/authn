package worker

import (
	"context"
	"github.com/imtanmoy/logx"
	"runtime"
)

type Dispatcher struct {
	maxWorkers int
	workerPool chan chan Job
	jobQueue   chan Job
	quit       chan bool
}

func NewDispatcher() *Dispatcher {
	processors := maxParallelism()
	pool := make(chan chan Job, processors)
	jq := make(chan Job)
	return &Dispatcher{workerPool: pool, maxWorkers: processors, jobQueue: jq, quit: make(chan bool)}
}

func (d *Dispatcher) Run(ctx context.Context) {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.workerPool)
		worker.Start(ctx)
	}

	go d.dispatch(ctx)
}

func (d *Dispatcher) dispatch(ctx context.Context) {
	for {

		select {
		case job := <-d.jobQueue:
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.workerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		case <-ctx.Done():
		case <-d.quit:
			logx.Info("closing dispatcher")
			return
		}
	}
}

// Stop signals the worker to stop listening for work requests.
func (d *Dispatcher) Stop() {
	d.quit <- true
}

func (d *Dispatcher) Send(fn func()) {
	payload := Payload{Fn: fn}
	work := Job{payload}
	d.jobQueue <- work
}

// MaxParallelism returns maximum number of go-routine should run
func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}
