package async

import "log"

// A worker performs task processing.
type worker struct {
	id         int
	taskQueue  chan Task
	workerPool WorkerPool
	quit       chan bool
}

// Creates a new worker instance.
func newWorker(id int, workerPool WorkerPool) *worker {
	return &worker{
		id:         id,
		workerPool: workerPool,
		taskQueue:  make(chan Task),
		quit:       make(chan bool),
	}
}

// Process tasks until a quit signal is received.
func (w *worker) start() {
	log.Printf("async: spinning up worker %d", w.id)
	go func(id int, workerPool WorkerPool, taskQueue chan Task, quit chan bool) {
		for {
			workerPool <- taskQueue // Indicate we're ready for a task
			select {
			case <-quit:
				log.Printf("async: quit signal in worker %d", id)
				return
			case task := <-taskQueue:
				if err := task.Process(); err != nil {
					log.Printf("async: unable to process task: %+v", err)
				}
			}
		}
	}(w.id, w.workerPool, w.taskQueue, w.quit)
}

// Stop processing tasks.
func (w *worker) stop() {
	go func() {
		w.quit <- true
	}()
}
