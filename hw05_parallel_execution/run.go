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

type breaker struct {
	sync.RWMutex
	flag bool
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
	var wgDone sync.WaitGroup
	breaker := breaker{flag: false}
	defer close(jobResult)

	wgDone.Add(numWorkers)
	for i := 1; i <= numWorkers; i++ {
		go worker(fmt.Sprint("Worker", i), &jobQueue, &jobResult, &wgDone, &breaker)
	}

	successCount := 0
	errorCount := 0
	errorStream := ""
	for successCount+errorCount < len(tasks) && errorCount < maxErrors {
		jr := <-jobResult
		jobDesc := fmt.Sprintf("\r\nTask[%d]{%s} by %s", jr.idx, getFuncName(tasks[jr.idx]), jr.worker)
		if jr.err != nil {
			errorCount++
			errorStream += jobDesc
		} else {
			successCount++
		}
	}

	breaker.Lock()
	breaker.flag = true
	breaker.Unlock()
	wgDone.Wait()

	if errorCount >= maxErrors {
		ErrErrorsLimitExceeded = errors.New("Errors limit [" + fmt.Sprint(maxErrors) + "] exceeded:" + errorStream)
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(name string, jobQueue *chan taskJob, jobResult *chan taskResult, wgDone *sync.WaitGroup, brk *breaker) {
	for {
		jq, ok := <-(*jobQueue)
		(*brk).RLock()
		flag := (*brk).flag
		(*brk).RUnlock()
		if !ok || flag {
			break
		}
		(*jobResult) <- taskResult{jq.idx, jq.task(), name}
	}
	wgDone.Done()
}

func getFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
