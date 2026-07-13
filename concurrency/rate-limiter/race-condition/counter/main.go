package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var syncMutex sync.Mutex
	counter := 0

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			syncMutex.Lock()
			counter++
			syncMutex.Unlock()
		}()
	}

	wg.Wait()

	fmt.Println("Expected:", 1000)
	fmt.Println("Actual:  ", counter)
}
