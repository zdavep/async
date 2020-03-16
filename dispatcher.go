package async

import (
	"log"
	"runtime"
)

// AutoSize forces the dispatcher to calculate worker pool size.
const AutoSize int = -1

// MinWorkers in the minimum worker pool size
const MinWorkers int = 2

// WorkerPool is a type alias for a channel of task channels (ie where workers receive tasks).
type WorkerPool = chan chan Task

// Dispatcher dispatches tasks to a pool of workers for processing.
type Dispatcher struct {
	workerPool WorkerPool
	size       int
	workers    []*worker
	quit       chan bool
}

// Determine the worker pool size based on the number of available CPUs.
func autoSizePool() (size int) {
	size = runtime.NumCPU()/2 + 1
	if size < MinWorkers {
		size = MinWorkers
	}
	return
}

// NewDispatcher creates a new dispatcher instance.
func NewDispatcher(size int) *Dispatcher {

	if size < MinWorkers {
		size = autoSizePool()
	}

	d := &Dispatcher{
		workerPool: make(WorkerPool, size),
		size:       size,
		workers:    make([]*worker, size),
		quit:       make(chan bool),
	}

	for i := 0; i < d.size; i++ {
		d.workers[i] = newWorker(i, d.workerPool)
	}

	return d
}

// Start spins up workers, then creates a go-routine that receives tasks and dispatches them to the worker pool.
func (d *Dispatcher) Start() {
	for _, worker := range d.workers {
		worker.start()
	}
	go func(workerPool WorkerPool, quit chan bool) {
		for {
			select {
			case <-quit:
				log.Println("async: quit signal in dispatcher")
				return
			case t := <-TaskQueue:
				go func(task Task) {
					worker := <-workerPool // Blocks until a worker is available.
					worker <- task
				}(t)
			}
		}
	}(d.workerPool, d.quit)
}

// Stop shuts down all workers in the pool.
func (d *Dispatcher) Stop() {
	for _, worker := range d.workers {
		worker.stop()
	}
	go func() {
		d.quit <- true
	}()
}
