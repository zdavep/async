package async

import "log"

// A worker performs task processing.
type worker struct {
	queue chan Task
	pool  chan chan Task
	quit  chan bool
}

// Creates a new worker instance.
func newWorker(pool chan chan Task) *worker {
	return &worker{
		pool:  pool,
		queue: make(chan Task),
		quit:  make(chan bool),
	}
}

// Process tasks until a quit signal is received.
func (w *worker) start() {
	go func() {
		for {
			w.pool <- w.queue // Indicate we're ready for a task
			select {
			case <-w.quit:
				return
			case task := <-w.queue:
				if err := task.Process(); err != nil {
					log.Printf("async: unable to process task: %+v", err)
				}
			}
		}
	}()
}

// Stop processing tasks.
func (w *worker) stop() {
	w.quit <- true
}
