package main

import (
	"context"
	"fmt"
	"sync"
)

func main() {
	ctx := context.Background()
	var wg sync.WaitGroup
	resultChan := make(chan int)
	jobChan := make(chan int)
	workerCount := 3
	worker := func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case job, ok := <-jobChan:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case resultChan <- job:
				}
			}
		}
	}
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker()
	}
	go func() {
		defer close(jobChan)
		for i := 0; i < 1000; i++ {
			select {
			case <-ctx.Done():
				return
			case jobChan <- i:
			}
		}
	}()
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	results := make([]int, 0)
	for v := range resultChan {
		results = append(results, v)
	}
	fmt.Println("Expected length:", 1000)
	fmt.Println("Actual length:  ", len(results))
}
