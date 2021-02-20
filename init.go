package timer

import (
	"log"
	"os"
	"sync"
)

func NewScheduler() *TaskScheduler {
	return &TaskScheduler{
		tasks:   new(sync.Map),
		running: new(sync.Map),
		add:     make(chan TaskInterface),
		stop:    make(chan struct{}),
		remove:  make(chan string),
		Logger:  log.New(os.Stdout, "[Control]: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

var TS *TaskScheduler

func init() {
	TS = NewScheduler()
}

func GetTaskScheduler() *TaskScheduler {
	return TS
}
