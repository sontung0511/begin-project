package main

import (
	"context"
	"fmt"
	"sync"
)

type Result struct {
	name  string
	value int
}

func main() {
	var wg sync.WaitGroup
	counts := map[string]int{}
	words := []string{
		"go", "java", "go", "python",
		"go", "java", "rust", "go",
	}
	ctx := context.Background()
	resultChan := make(chan string)
	jobChan := make(chan string)
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
	wg.Add(len(words))
	for i := 0; i < len(words); i++ {
		go worker()
	}
	go func() {
		defer close(jobChan)
		for _, word := range words {
			select {
			case <-ctx.Done():
				return
			case jobChan <- word:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()
	for result := range resultChan {
		counts[result]++
	}
	fmt.Println(counts)
}
