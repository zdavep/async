package async

import (
	"log"
	"runtime"
)

// AutoSize forces the dispatcher to calculate worker pool size.
const AutoSize int = -1

// MinWorkers in the minimum worker pool size
const MinWorkers int = 2

// Dispatcher dispatches tasks to a pool of workers for processing.
type Dispatcher struct {
	pool    workerPool
	size    int
	workers []*worker
	quit    chan bool
	errs    chan error
}

// Determine the worker pool size based on the number of available CPUs.
func autoSizePool() (size int) {
	size = runtime.NumCPU()/2 + 1
	if size < MinWorkers {
		size = MinWorkers
	}
	return
}

// NewDispatcher creates a new dispatcher and returns it with an error channel.
func NewDispatcher(size int) (*Dispatcher, <-chan error) {
	if size < MinWorkers {
		size = autoSizePool()
	}
	errs := make(chan error)
	d := &Dispatcher{
		pool:    make(workerPool, size),
		size:    size,
		workers: make([]*worker, size),
		quit:    make(chan bool),
		errs:    errs,
	}
	for i := 0; i < d.size; i++ {
		d.workers[i] = newWorker(i, d.pool, errs)
	}
	return d, errs
}

// LogErrors is a helper for logging errors to stderr.
func (d *Dispatcher) LogErrors(errs <-chan error) {
	go func() {
		for err := range errs {
			log.Printf("async: process error: %v", err)
		}
	}()
}

// Start spins up workers, then creates a go-routine that dispatches tasks to the worker pool.
func (d *Dispatcher) Start() {
	for _, worker := range d.workers {
		worker.start()
	}
	go dispatchTasks(d.pool, d.quit)
}

// Dispatch tasks until instructed to quit.
func dispatchTasks(pool workerPool, quit chan bool) {
	for {
		select {
		case <-quit:
			return
		case t := <-TaskQueue:
			go func(task Task) {
				worker := <-pool // Blocks until a worker is available.
				worker <- task
			}(t)
		}
	}
}

// Stop shuts down all workers in the pool.
func (d *Dispatcher) Stop() {
	defer close(d.errs)
	for _, worker := range d.workers {
		worker.stop()
	}
	go func() {
		d.quit <- true
	}()
}
