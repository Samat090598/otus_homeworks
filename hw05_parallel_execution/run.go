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
		tasksCh = make(chan Task)
		errCh   = make(chan error, m)
		stopCh  = make(chan struct{})
	)

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for {
				select {
				case task, ok := <-tasksCh:
					if !ok {
						return
					}

					if err := task(); err != nil {
						errCh <- err
					}
				case <-stopCh:
					return
				}
			}
		}()
	}

	go func() {
		for _, task := range tasks {
			select {
			case tasksCh <- task:

			case <-stopCh:
				close(tasksCh)
				return
			}
		}
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
