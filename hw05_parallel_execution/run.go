package hw05parallelexecution

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

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
	return RunLogged(tasks, numWorkers, maxErrors, nil)
}

func RunLogged(tasks []Task, numWorkers, maxErrors int, out io.StringWriter) error {
	var workerBufferSize int
	if numWorkers <= 1 {
		numWorkers = 1       // Хотя бы один worker должен быть запущен
		workerBufferSize = 0 // Длинный буфер не имеет смысла
	} else {
		workerBufferSize = len(tasks) - numWorkers + 1 // Если все короткие задачи упадут на 1 worker
		if workerBufferSize < 0 {
			workerBufferSize = 0
		}
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
		jobResults[i] = worker(fmt.Sprint("Worker", i+1), jobQueue, &wgDone, breaker, workerBufferSize)
	}

	successCount := 0
	errorCount := 0
	for successCount+errorCount < len(tasks) && errorCount < maxErrors {
		jr := getNextWorkerResult(jobResults)
		jobMess := fmt.Sprintf("\r\nTask[%d]{%s} by %s", jr.idx, getFuncName(tasks[jr.idx]), jr.worker)
		if jr.err != nil {
			errorCount++
			jobMess += " failed"
		} else {
			successCount++
			jobMess += " completed"
		}
		if out != nil {
			out.WriteString(jobMess)
		}
	}

	close(breaker)
	wgDone.Wait()

	if errorCount >= maxErrors {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(name string, jobQueue <-chan taskJob,
	wgDone *sync.WaitGroup, brk <-chan struct{},
	bufferSize int,
) <-chan taskResult {
	var jq taskJob
	var jr taskResult
	jobResult := make(chan taskResult, bufferSize)
	go func() {
		defer close(jobResult)
		defer wgDone.Done()
		jget := false
		for {
			select {
			case <-brk:
				return
			case jq, jget = <-jobQueue:
			}
			if jget {
				jr = taskResult{jq.idx, jq.task(), name}
				select {
				case <-brk:
					return
				case jobResult <- jr:
				}
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
