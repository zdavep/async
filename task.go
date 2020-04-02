package async

import (
	"log"
	"os"
	"strconv"
	"strings"
)

// The default queue buffer size
const defSize int = 1000

// Task is a unit of work that needs to be processed.
type Task interface {
	Process() error
}

// TaskQueue is a buffered channel for processing tasks.
var TaskQueue chan Task

// Initialize the task queue.
func init() {
	var size int
	var err error
	envSize := strings.TrimSpace(os.Getenv("ASYNC_TASK_QUEUE_SIZE"))
	if envSize == "" {
		size = defSize
	} else if size, err = strconv.Atoi(envSize); err != nil {
		log.Panicf("async: unable to convert task queue size to int: %v", err)
	}
	TaskQueue = make(chan Task, size)
}
