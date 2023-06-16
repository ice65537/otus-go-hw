package memstore

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestStorage(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("Create event", func(t *testing.T) {
		store := New()
		ctx := context.Background()
		log := logger.New("test", "DEBUG", 5)
		store.Init(ctx, log)
		t1 := time.Now()
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("Too much workers", func(t *testing.T) {
		tasksCount := 30
		tasks := make([]Task, 0, tasksCount)

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

		workersCount := 300
		maxErrorsCount := 20

		err := RunLogged(tasks, workersCount, maxErrorsCount, os.Stdout)
		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})

	t.Run("Zero (one) worker and Nonsequental and skip all errors tests", func(t *testing.T) {
		tasksCount := 100
		tasks := make([]Task, 0, tasksCount)

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
		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed at 1st run")

		runTasksCount = 0
		workersCount = 50
		maxErrorsCount = -1 // ignore errors
		start = time.Now()
		err = Run(tasks, workersCount, maxErrorsCount)
		elapsedTime2 := time.Since(start)
		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed at 2nd run")
		require.LessOrEqual(t, int64(elapsedTime2), int64(elapsedTime), "tasks were run sequentially?")
	})
}
