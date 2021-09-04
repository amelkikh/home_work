package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrInvalidWorkerCount  = errors.New("negative worker count")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workerCount, errorLimit int) error {
	if workerCount <= 0 {
		return ErrInvalidWorkerCount
	}
	if errorLimit <= 0 {
		errorLimit = len(tasks) + 1
	}

	var wg sync.WaitGroup
	var errCnt int32
	taskCh := make(chan Task)
	wg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()
			for task := range taskCh {
				if err := task(); err != nil {
					atomic.AddInt32(&errCnt, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt32(&errCnt) >= int32(errorLimit) {
			break
		}
		taskCh <- task
	}

	close(taskCh)
	wg.Wait()

	if errCnt >= int32(errorLimit) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
