package worker

import (
	"context"
	"github.com/imtanmoy/logx"
)

type Payload struct {
	Fn func()
}

// Job represents the job to be run
type Job struct {
	Payload Payload
}

// Worker represents the worker that executes the job
type Worker struct {
	workerPool chan chan Job
	jobChannel chan Job
	quit       chan bool
}

func NewWorker(workerPool chan chan Job) *Worker {
	return &Worker{
		workerPool: workerPool,
		jobChannel: make(chan Job),
		quit:       make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w *Worker) Start(ctx context.Context) {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.workerPool <- w.jobChannel

			select {
			case job := <-w.jobChannel:
				// we have received a work request.
				job.Payload.Fn()
			case <-w.quit:
			case <-ctx.Done():
				// we have received a signal to stop
				logx.Info("closing worker...")
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w *Worker) Stop() {
	w.quit <- true
}
