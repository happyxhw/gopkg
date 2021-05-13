package dispatcher

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrStopped     = errors.New("dispatcher stopped")
	ErrStopTimeout = errors.New("stop timeout")
	ErrJobTimeout  = errors.New("job timeout")
)

// Job task job
type Job func() (interface{}, error)

type Task struct {
	Job     Job
	Timeout time.Duration
}

// ResultQueue job result chan
type ResultQueue chan *JobResult

// JobResult job result
type JobResult struct {
	Result interface{}
	Error  error
}

type worker struct {
	stopChan chan struct{}
	taskChan chan *Task
}

func newWorker() *worker {
	return &worker{
		stopChan: make(chan struct{}),
		taskChan: make(chan *Task),
	}
}

// start worker & waiting for task
func (w *worker) start(pool chan *worker, resultCh ResultQueue) {
	go func() {
		for {
			pool <- w
			task, ok := <-w.taskChan
			if !ok {
				w.stopChan <- struct{}{}
				return
			}
			var jobRes JobResult
			if task.Timeout <= 0 {
				res, err := task.Job()
				jobRes.Result = res
				jobRes.Error = err
				resultCh <- &jobRes
				return
			}
			resCh := make(chan interface{})
			errCh := make(chan error)
			ctx, cancel := context.WithTimeout(context.TODO(), task.Timeout)
			go func() {
				res, err := task.Job()
				if err != nil {
					errCh <- err
					return
				}
				resCh <- res
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
			resultCh <- &jobRes
		}
	}()
}

type Options struct {
	PoolSize    int
	StopTimeout time.Duration
}

var defaultOpts = &Options{
	PoolSize:    1000,
	StopTimeout: time.Second * 100,
}

type Option func(*Options)

func WithPoolSize(size int) Option {
	return func(opts *Options) {
		opts.PoolSize = size
	}
}

func WithStopTimeout(timeout time.Duration) Option {
	return func(opts *Options) {
		opts.StopTimeout = timeout
	}
}

type Dispatcher struct {
	sync.RWMutex

	pool        chan *worker
	taskQueue   chan *Task
	resultQueue ResultQueue
	stopChan    chan struct{}
	stopTimeout time.Duration
	stopped     bool
}

// NewDispatcher init Dispatcher
func NewDispatcher(opts ...Option) *Dispatcher {
	options := defaultOpts
	for _, o := range opts {
		o(options)
	}
	d := &Dispatcher{
		taskQueue:   make(chan *Task),
		stopChan:    make(chan struct{}),
		resultQueue: make(ResultQueue),
		pool:        make(chan *worker, options.PoolSize),
		stopTimeout: options.StopTimeout,
	}
	for i := 0; i < options.PoolSize; i++ {
		worker := newWorker()
		worker.start(d.pool, d.resultQueue)
	}
	go d.dispatch()
	return d
}

// start Dispatcher
func (d *Dispatcher) dispatch() {
	for {
		task, ok := <-d.taskQueue
		if ok {
			worker := <-d.pool
			worker.taskChan <- task
			continue
		}
		d.stopChan <- struct{}{}
		return
	}
}

// Stop Dispatcher
func (d *Dispatcher) Stop() error {
	d.Lock()
	d.stopped = true
	close(d.taskQueue)
	d.Unlock()
	<-d.stopChan
	doneCh := make(chan struct{})
	ctx, cancel := context.WithTimeout(context.TODO(), d.stopTimeout)
	defer cancel()
	go func() {
		for i := 0; i < cap(d.pool); i++ {
			worker := <-d.pool
			close(worker.taskChan)
			<-worker.stopChan
		}
		doneCh <- struct{}{}
	}()
	var err error
	select {
	case <-doneCh:
	case <-ctx.Done():
		err = ErrStopTimeout
	}
	close(d.resultQueue)
	return err
}

// Send task to Dispatcher
func (d *Dispatcher) Send(task *Task) error {
	d.RLock()
	defer d.RUnlock()
	if d.stopped {
		return ErrStopped
	}
	d.taskQueue <- task
	return nil
}

// ResultCh get the result chan
func (d *Dispatcher) ResultCh() ResultQueue {
	return d.resultQueue
}
