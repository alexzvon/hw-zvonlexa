package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var errCount int
	ct := make(chan Task)
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for task := range ct {
				mu.Lock()
				errCountGo := errCount
				mu.Unlock()

				switch {
				case m == 0:
					_ = task()
				case errCountGo < m:
					if err := task(); err != nil {
						mu.Lock()
						errCount++
						mu.Unlock()
					}
				}
			}
		}()
	}

	go func() {
		defer close(ct)

		for _, task := range tasks {
			mu.Lock()
			errCountGo := errCount
			mu.Unlock()

			if errCountGo >= m {
				break
			}

			ct <- task
		}
	}()

	wg.Wait()

	if errCount < m {
		return nil
	}

	return ErrErrorsLimitExceeded
}
