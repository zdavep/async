package async

import "log"

// A worker performs task processing.
type worker struct {
	id    int
	queue chan Task
	pool  chan chan Task
	quit  chan bool
}

// Creates a new worker instance.
func newWorker(id int, pool chan chan Task) *worker {
	return &worker{
		id:    id,
		pool:  pool,
		queue: make(chan Task),
		quit:  make(chan bool),
	}
}

// Process tasks until a quit signal is received.
func (w *worker) start() {
	go func(pool chan chan Task, queue chan Task, quit chan bool) {
		for {
			pool <- queue // Indicate we're ready for a task
			select {
			case <-quit:
				log.Printf("async: quit signal in worker %d", w.id)
				return
			case task := <-queue:
				if err := task.Process(); err != nil {
					log.Printf("async: unable to process task: %+v", err)
				}
			}
		}
	}(w.pool, w.queue, w.quit)
}

// Stop processing tasks.
func (w *worker) stop() {
	go func() {
		w.quit <- true
	}()
}
