package main

import (
	"context"
	"fmt"
	"sync"
)

type Result struct {
	Index int
	Value int
}

func ProcessOrders(ctx context.Context, orders []int, workCount int) ([]int, error) {
	result := make([]int, len(orders))
	if len(orders) == 0 {
		return result, nil
	}
	if workCount <= 0 {
		return nil, fmt.Errorf("work count is zero")
	}
	if workCount > len(orders) {
		workCount = len(orders)
	}
	resultChan := make(chan Result)
	jobChan := make(chan Result)
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
				order := job.Value + job.Value/10
				select {
				case <-ctx.Done():
					return
				case resultChan <- Result{Index: job.Index, Value: order}:
				}
			}
		}
	}
	wg.Add(workCount)
	for i := 0; i < workCount; i++ {
		go worker()
	}
	go func() {
		defer close(jobChan)
		for index, order := range orders {
			select {
			case <-ctx.Done():
				return
			case jobChan <- Result{Index: index, Value: order}:
			}
		}
	}()
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	for order := range resultChan {
		result[order.Index] = order.Value
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
