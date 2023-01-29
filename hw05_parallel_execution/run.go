package hw05parallelexecution

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
)

var ErrErrorsLimitExceeded error

type Task func() error

type taskJob struct {
	idx  int
	task Task
}

type taskResult struct {
	idx    int
	err    error
	worker string
}

// var taskQueue chan Task

// Run starts tasks in numWorkers goroutines and stops its work when receiving maxErrors errors from tasks.
func Run(tasks []Task, numWorkers, maxErrors int) error {
	if numWorkers <= 1 {
		numWorkers = 1 // Хотя бы один worker должен быть запущен
	}
	if maxErrors <= 0 {
		maxErrors = len(tasks) + 100 // Обработка с игнорированием ошибок
	}

	jobQueue := make(chan taskJob, len(tasks))
	for i, task := range tasks {
		jobQueue <- taskJob{i, task}
	}
	close(jobQueue)

	jobResult := make(chan taskResult, len(tasks))
	massBreak := make(chan struct{})
	breakFlag := false
	defer close(jobResult)
	defer close(massBreak)

	for i := 1; i <= numWorkers; i++ {
		go worker(fmt.Sprint("Worker", i), &jobQueue, &jobResult, &massBreak, &breakFlag)
	}

	successCount := 0
	errorCount := 0
	errorStream := ""
	for successCount+errorCount < len(tasks) && errorCount < maxErrors {
		jr := <-jobResult
		jobDesc := fmt.Sprintf("\r\nTask[%d]{%s} by %s", jr.idx, runtime.FuncForPC(reflect.ValueOf(tasks[jr.idx]).Pointer()).Name(), jr.worker)
		if jr.err != nil {
			errorCount++
			errorStream += jobDesc
		} else {
			successCount++
		}
	}

	breakFlag = true
	for i := 1; i <= numWorkers; i++ {
		massBreak <- struct{}{}
	}

	if errorCount >= maxErrors {
		ErrErrorsLimitExceeded = errors.New("Errors limit [" + fmt.Sprint(maxErrors) + "] exceeded:" + errorStream)
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(name string, jobQueue *chan taskJob, jobResult *chan taskResult, breaker *chan struct{}, breakFlag *bool) {
	for {
		select {
		case <-(*breaker):
			return
		case jq, ok := <-(*jobQueue):
			if ok && !*breakFlag {
				(*jobResult) <- taskResult{jq.idx, jq.task(), name}
			}
		}
	}
}
