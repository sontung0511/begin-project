package main

import (
	"context"
	"fmt"
	"sync"
)

// Problem: Process Orders with Worker Pool
//
// You are given a list of order amounts and a fixed number of workers.
//
// Each worker receives an order through a shared job channel,
// applies a 10% service fee, and sends the processed result
// to a shared result channel.
//
// Requirements:
//
//  1. Create exactly workerCount worker goroutines.
//
//  2. Send all orders through a shared job channel.
//
//  3. Each worker calculates:
//
//     finalAmount = orderAmount + orderAmount/10
//
//  4. Preserve the original order of the input orders.
//
//  5. Support context cancellation.
//
//  6. Return an error when workerCount <= 0.
//
//  7. Do not use time.Sleep for synchronization.
//
//  8. Do not leak goroutines.
//
// Function:
//
//	func ProcessOrders(
//		ctx context.Context,
//		orders []int,
//		workerCount int,
//	) ([]int, error)
//
// Example:
//
//	Input:
//		orders      = []int{100, 200, 300, 400}
//		workerCount = 2
//
//	Expected:
//		[]int{110, 220, 330, 440}
type Job struct {
	index int
	value int
}
type Result struct {
	index int
	value int
}

func ProcessOrders(
	ctx context.Context,
	orders []int,
	workerCount int,
) ([]int, error) {
	if orders == nil || len(orders) == 0 {
		return []int{}, nil
	}
	if workerCount < 1 {
		return nil, fmt.Errorf("workerCount must be greater than zero")
	}
	if workerCount > len(orders) {
		workerCount = len(orders)
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
				orderAmount := job.value + job.value/10
				select {
				case <-ctx.Done():
					return
				case resultChan <- Result{index: job.index, value: orderAmount}:
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
		for index, order := range orders {
			select {
			case <-ctx.Done():
				return
			case jobChan <- Job{index, order}:
			}
		}
	}()
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	results := make([]int, len(orders))
	received := 0
	for received < len(orders) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result, ok := <-resultChan:
			if !ok {
				return nil, fmt.Errorf("channel closed")
			}
			results[result.index] = result.value
			received++
		}
	}

	return results, nil
}
