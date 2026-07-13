package main

import (
	"context"
	"sync"
)

// Package workerpool implements concurrent job processing using the Worker Pool pattern.
//
// Problem: Concurrent Job Processing with Worker Pool
//
// You are given a list of integer jobs and a fixed number of workers.
//
// Each worker must run in its own goroutine and receive jobs through a channel.
// For each job, the worker calculates the square of the input value.
//
// Requirements:
//
//  1. Create exactly workerCount worker goroutines.
//  2. Send all jobs to workers through a channel.
//  3. Preserve the original order of the input jobs in the returned result.
//  4. Stop processing when the context is cancelled.
//  5. Return ctx.Err() when cancellation occurs.
//  6. Return an error when workerCount is less than or equal to zero.
//  7. Do not use time.Sleep for synchronization.
//  8. Do not leave any goroutine running after the function returns.
//
// Function:
//
//	func ProcessJobs(
//		ctx context.Context,
//		jobs []int,
//		workerCount int,
//	) ([]int, error)
//
// Example:
//
//	Input:
//		jobs       = []int{1, 2, 3, 4, 5}
//		workerCount = 3
//
//	Expected:
//		[]int{1, 4, 9, 16, 25}
//
// Even though jobs may finish in a different order, the returned result must
// match the original input order.
type Job struct {
	index int
	value int
}
type Result struct {
	index int
	value int
}

func ProcessJobs(ctx context.Context, jobs []int, workerCount int) ([]int, error) {
	if len(jobs) == 0 || workerCount == 0 {
		return make([]int, 0), nil
	}
	if workerCount > len(jobs) {
		workerCount = len(jobs)
	}
	resultChan := make(chan Result)
	jobChan := make(chan Job)
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
				result := Result{
					index: job.index,
					value: job.value * job.value,
				}
				select {
				case <-ctx.Done():
					return
				case resultChan <- result:
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
		for index, value := range jobs {
			select {
			case <-ctx.Done():
				return
			case jobChan <- Job{index, value}:
			}
		}
	}()
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	results := make([]int, len(jobs))
	received := 0
	for received < len(jobs) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result, ok := <-resultChan:
			if !ok {
				return nil, ctx.Err()
			}
			results[result.index] = result.value
			received++
		}
	}
	return results, nil
}
