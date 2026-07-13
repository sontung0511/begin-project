package main

import (
	"context"
	"fmt"
	"sync"
)

// Problem: Sum Squares with Fan-Out / Fan-In
//
// You are given a list of integers and a fixed number of workers.
//
// Split the work across multiple goroutines.
// Each worker receives numbers from a shared input channel,
// calculates the square of each number, and sends the result
// to a shared output channel.
//
// Requirements:
//
//  1. Create exactly workerCount worker goroutines.
//  2. Send input values through a channel.
//  3. Each worker calculates value * value.
//  4. Merge all worker results into one output channel.
//  5. Return the total sum of all squared values.
//  6. Support context cancellation.
//  7. Return an error when workerCount <= 0.
//  8. Do not use time.Sleep for synchronization.
//  9. Do not leak goroutines.
//
// Function:
//
//	func SumSquares(
//		ctx context.Context,
//		numbers []int,
//		workerCount int,
//	) (int, error)
//
// Example:
//
//	Input:
//		numbers     = []int{1, 2, 3, 4, 5}
//		workerCount = 3
//
//	Expected:
//		55
//
// Explanation:
//
//	1² + 2² + 3² + 4² + 5²
//	= 1 + 4 + 9 + 16 + 25
//	= 55
func SumSquares(
	ctx context.Context,
	numbers []int,
	workerCount int,
) (int, error) {
	if len(numbers) == 0 {
		return 0, fmt.Errorf("no numbers provided")
	}
	if workerCount <= 0 {
		return 0, fmt.Errorf("invalid number of workers")
	}
	if workerCount > len(numbers) {
		workerCount = len(numbers)
	}
	resultChan := make(chan int)
	jobChan := make(chan int)
	var wg sync.WaitGroup
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
				square := job * job

				select {
				case <-ctx.Done():
					return
				case resultChan <- square:
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
		for _, number := range numbers {
			select {
			case <-ctx.Done():
				return
			case jobChan <- number:

			}

		}
	}()
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	total := 0
	received := 0
	for received < len(numbers) {
		select {
		case <-ctx.Done():
			return total, ctx.Err()
		case result, ok := <-resultChan:
			if !ok {
				return total, ctx.Err()
			}
			total += result
			received++
		}
	}
	return total, nil
}
