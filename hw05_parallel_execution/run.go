package hw05parallelexecution

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sync"
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

	jobResults := make([]<-chan taskResult, numWorkers)
	var wgDone sync.WaitGroup
	breaker := make(chan struct{})

	wgDone.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		jobResults[i] = worker(fmt.Sprint("Worker", i+1), jobQueue, &wgDone, breaker, len(tasks))
	}

	successCount := 0
	errorCount := 0
	errorStream := ""
	for successCount+errorCount < len(tasks) && errorCount < maxErrors {
		jr := getNextWorkerResult(jobResults)
		jobDesc := fmt.Sprintf("\r\nTask[%d]{%s} by %s", jr.idx, getFuncName(tasks[jr.idx]), jr.worker)
		if jr.err != nil {
			errorCount++
			errorStream += jobDesc
		} else {
			successCount++
		}
	}

	close(breaker)
	wgDone.Wait()

	if errorCount >= maxErrors {
		ErrErrorsLimitExceeded = errors.New("Errors limit [" + fmt.Sprint(maxErrors) + "] exceeded:" + errorStream)
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(name string, jobQueue <-chan taskJob, wgDone *sync.WaitGroup,
	brk <-chan struct{}, bufferSize int,
) <-chan taskResult {
	jobResult := make(chan taskResult, bufferSize)
	go func() {
		defer close(jobResult)
		defer wgDone.Done()
		for {
			select {
			case jq, ok := <-jobQueue:
				if ok {
					jobResult <- taskResult{jq.idx, jq.task(), name}
				}
			case <-brk:
				return
			}
		}
	}()
	return jobResult
}

func getFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func getNextWorkerResult(jobResults []<-chan taskResult) taskResult {
	for {
		for _, ch := range jobResults {
			select {
			case jobResult := <-ch:
				return jobResult
			default:
			}
		}
	}
}
