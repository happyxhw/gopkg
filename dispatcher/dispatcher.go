package dispatcher

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrStopped    = errors.New("dispatcher stopped")
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

type dispatcher struct {
	pool        chan *worker
	jobQueue    chan Job
	resultQueue ResultQueue
	stopChan    chan struct{}
	stopped     bool
	lock        *sync.RWMutex
}

// NewDispatcher init dispatcher
func NewDispatcher(workerSize int, jobTimeout time.Duration) *dispatcher {
	d := &dispatcher{
		jobQueue:    make(chan Job),
		stopChan:    make(chan struct{}),
		resultQueue: make(ResultQueue),
		pool:        make(chan *worker, workerSize),
		lock:        &sync.RWMutex{},
	}

	for i := 0; i < workerSize; i++ {
		worker := newWorker(jobTimeout)
		worker.start(d.pool, d.resultQueue)
	}
	go d.dispatcher()
	return d
}

// start dispatcher
func (d *dispatcher) dispatcher() {
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

// Stop dispatcher
func (d *dispatcher) Stop() {
	d.lock.Lock()
	d.stopped = true
	d.lock.Unlock()
	close(d.jobQueue)
	<-d.stopChan

	for i := 0; i < cap(d.pool); i++ {
		worker := <-d.pool
		close(worker.jobChan)
		<-worker.stopChan
	}
	close(d.resultQueue)
}

// Send task to dispatcher
func (d *dispatcher) Send(job Job) error {
	d.lock.RLock()
	defer d.lock.RUnlock()
	if d.stopped {
		return ErrStopped
	}
	d.jobQueue <- job
	return nil
}

// ResultCh: get the result chan
func (d *dispatcher) ResultCh() ResultQueue {
	return d.resultQueue
}
