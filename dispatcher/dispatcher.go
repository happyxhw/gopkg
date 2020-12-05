package dispatcher

import (
	"context"
	"errors"
	"time"

	"go.uber.org/atomic"
)

var (
	ErrStopped    = errors.New("Dispatcher stopped")
	ErrJobTimeout = errors.New("job timeout")
)

// Job task job
type Job func() (interface{}, error)

// ResultQueue job result chan
type ResultQueue chan *JobResult

// JobResult job result
type JobResult struct {
	Result interface{}
	Error  error
}

type worker struct {
	stopChan   chan struct{}
	jobChan    chan Job
	jobTimeout time.Duration
}

func newWorker(jobTimeout time.Duration) *worker {
	return &worker{
		stopChan:   make(chan struct{}),
		jobChan:    make(chan Job),
		jobTimeout: jobTimeout,
	}
}

// start worker & waiting for task
func (w *worker) start(pool chan *worker, resultCh ResultQueue) {
	go func() {
		for {
			pool <- w
			job, ok := <-w.jobChan
			if ok {
				var jobRes JobResult
				if w.jobTimeout <= 0 {
					res, err := job()
					jobRes.Result = res
					jobRes.Error = err
				} else {
					resCh := make(chan interface{})
					errCh := make(chan error)
					ctx, cancel := context.WithTimeout(context.TODO(), w.jobTimeout)
					go func() {
						res, err := job()
						if err != nil {
							errCh <- err
							close(errCh)
						} else {
							resCh <- res
							close(resCh)
						}
					}()
					select {
					case res := <-resCh:
						jobRes.Result = res
					case err := <-errCh:
						jobRes.Error = err
					case <-ctx.Done():
						jobRes.Error = ErrJobTimeout
					}
					cancel()
				}
				resultCh <- &jobRes
			} else {
				w.stopChan <- struct{}{}
				return
			}
		}
	}()
}

type Dispatcher struct {
	pool        chan *worker
	jobQueue    chan Job
	resultQueue ResultQueue
	stopChan    chan struct{}
	stopped     atomic.Bool
}

// NewDispatcher init Dispatcher
func NewDispatcher(workerSize int, jobTimeout time.Duration) *Dispatcher {
	d := &Dispatcher{
		jobQueue:    make(chan Job),
		stopChan:    make(chan struct{}),
		resultQueue: make(ResultQueue),
		pool:        make(chan *worker, workerSize),
	}

	for i := 0; i < workerSize; i++ {
		worker := newWorker(jobTimeout)
		worker.start(d.pool, d.resultQueue)
	}
	go d.dispatcher()
	return d
}

// start Dispatcher
func (d *Dispatcher) dispatcher() {
	for {
		job, ok := <-d.jobQueue
		if ok {
			worker := <-d.pool
			worker.jobChan <- job
		} else {
			d.stopChan <- struct{}{}
			return
		}
	}
}

// Stop Dispatcher
func (d *Dispatcher) Stop() {
	d.stopped.Swap(true)
	close(d.jobQueue)
	<-d.stopChan

	for i := 0; i < cap(d.pool); i++ {
		worker := <-d.pool
		close(worker.jobChan)
		<-worker.stopChan
	}
	close(d.resultQueue)
}

// Send task to Dispatcher
func (d *Dispatcher) Send(job Job) error {
	if d.stopped.Load() {
		return ErrStopped
	}
	d.jobQueue <- job
	return nil
}

// ResultCh: get the result chan
func (d *Dispatcher) ResultCh() ResultQueue {
	return d.resultQueue
}
