package main

import (
	"errors"
	"os"
	"sync/atomic"
	"time"

	hwpe "github.com/ice65537/otus-go-hw/hw05_parallel_execution"
)

func main() {
	tasksCount := 100
	tasks := make([]hwpe.Task, 0, tasksCount)

	var runTasksCount int32

	for i := 0; i < tasksCount; i++ {
		if i%2 == 0 {
			tasks = append(tasks, func() error {
				time.Sleep(50 * time.Millisecond)
				atomic.AddInt32(&runTasksCount, 1)
				return errors.New("Err")
			})
		} else {
			tasks = append(tasks, func() error {
				time.Sleep(50 * time.Millisecond)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}
	}

	runTasksCount = 0
	workersCount := -45     // equal to 1
	maxErrorsCount := -1345 // ignore errors
	hwpe.RunLogged(tasks, workersCount, maxErrorsCount, os.Stdout)

	runTasksCount = 0
	workersCount = 50
	maxErrorsCount = -1 // ignore errors
	hwpe.RunLogged(tasks, workersCount, maxErrorsCount, os.Stdout)
}
