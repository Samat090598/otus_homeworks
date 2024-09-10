package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n int, m int) error {
	var (
		wg      sync.WaitGroup
		tasksCh = make(chan Task, len(tasks))
		errCh   = make(chan struct{}, m)
		stopCh  = make(chan struct{})
	)

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for task := range tasksCh {
				select {
				case <-stopCh:
					return
				default:
					err := task()
					if err != nil {
						errCh <- struct{}{}
					}
				}
			}
		}()
	}

	go func() {
		for _, task := range tasks {
			tasksCh <- task
		}
		close(tasksCh)
	}()

	errCount := 0
	go func() {
		for range errCh {
			errCount++
			if errCount >= m {
				stopCh <- struct{}{}
				close(stopCh)
				return
			}
		}
	}()

	wg.Wait()
	close(errCh)

	if errCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
