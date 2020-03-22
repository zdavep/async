package async

import "log"

// A type alias for a task channel (ie where workers receive tasks).
type workQueue = chan Task

// A type alias for a channel of work queues.
type workerPool = chan workQueue

// A worker performs task processing.
type worker struct {
	id    int
	queue workQueue  // Where tasks are received.
	pool  workerPool // Where this worker asks for a task.
	quit  chan bool
}

// Creates a new worker instance.
func newWorker(id int, pool workerPool) *worker {
	return &worker{
		id:    id,
		pool:  pool,
		queue: make(workQueue),
		quit:  make(chan bool),
	}
}

// Start task processing go-routine.
func (w *worker) start() {
	go processTasks(w.id, w.pool, w.queue, w.quit)
}

// Process tasks until instructed to stop.
func processTasks(id int, pool workerPool, worker workQueue, quit chan bool) {
	for {
		pool <- worker // Indicate worker is ready for a task
		select {
		case <-quit:
			return
		case task := <-worker:
			if err := task.Process(); err != nil {
				log.Printf("async: unable to process task: %+v", err)
			}
		}
	}
}

// Stop processing tasks.
func (w *worker) stop() {
	go func() {
		w.quit <- true
	}()
}
