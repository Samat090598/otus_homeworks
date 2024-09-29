package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n int, m int) error {
	var (
		wg       sync.WaitGroup
		errCount int64
		taskChan = make(chan Task, len(tasks))
	)

	if m < 1 {
		m = 1
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				if int(atomic.LoadInt64(&errCount)) >= m {
					return
				}

				err := task()
				if err != nil {
					atomic.AddInt64(&errCount, 1)
				}
			}
		}()
	}

	go func() {
		for _, task := range tasks {
			taskChan <- task
		}
		close(taskChan)
	}()

	wg.Wait()

	if int(atomic.LoadInt64(&errCount)) >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
