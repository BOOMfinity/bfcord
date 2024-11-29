package utils

import "sync"

func RunParallel(tasks ...func()) {
	wg := new(sync.WaitGroup)
	wg.Add(len(tasks))
	for _, task := range tasks {
		go func(t func()) {
			t()
			wg.Done()
		}(task)
	}
	wg.Wait()
}
