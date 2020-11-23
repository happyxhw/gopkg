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

type Job func() (interface{}, error)
type ResultQueue chan *JobResult

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

func (w *worker) start(pool chan *worker, resultCh ResultQueue) {
	go func() {
		for {
			pool <- w
			job, ok := <-w.jobChan
			if ok {
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
				var jobRes JobResult
				select {
				case res := <-resCh:
					jobRes.Result = res
				case err := <-errCh:
					jobRes.Error = err
				case <-ctx.Done():
					jobRes.Error = ErrJobTimeout
				}
				resultCh <- &jobRes
				cancel()
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
	size        int
	stopChan    chan struct{}
	stopped     bool
	lock        *sync.RWMutex
}

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

func (d *dispatcher) Send(job Job) error {
	d.lock.RLock()
	defer d.lock.RUnlock()
	if d.stopped {
		return ErrStopped
	}
	d.jobQueue <- job
	return nil
}

func (d *dispatcher) ResultCh() ResultQueue {
	return d.resultQueue
}
