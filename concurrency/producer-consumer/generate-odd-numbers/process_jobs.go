package main

import (
	"context"
	"fmt"
	"sync"
)

func GenerateOddNumbers(
	ctx context.Context,
	producerCount int,
	consumerCount int,
	itemsPerProducer int,
) ([]int, error) {
	if producerCount < 1 {
		return nil, fmt.Errorf("producerCount must be greater than zero")
	}
	if consumerCount < 1 {
		return nil, fmt.Errorf("consumerCount must be greater than zero")
	}
	if itemsPerProducer < 1 {
		return nil, fmt.Errorf("itemsPerProducer must be greater than zero")
	}
	resultChan := make(chan int, consumerCount)
	jobChan := make(chan int, producerCount)
	var producerWG sync.WaitGroup
	var consumerWG sync.WaitGroup
	producer := func(producerID int) {
		defer producerWG.Done()
		for i := 0; i < itemsPerProducer; i++ {
			start := producerID*itemsPerProducer + i
			order := 1 + start*2
			select {
			case <-ctx.Done():
				return
			case jobChan <- order:
			}
		}
	}
	consumer := func() {
		defer consumerWG.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case job, ok := <-jobChan:
				if !ok {
					return
				}
				step := job
				select {
				case <-ctx.Done():
					return
				case resultChan <- step:
				}
			}
		}
	}
	producerWG.Add(producerCount)
	for i := 0; i < producerCount; i++ {
		go producer(i)
	}
	consumerWG.Add(consumerCount)
	for i := 0; i < consumerCount; i++ {
		go consumer()
	}
	go func() {
		producerWG.Wait()
		close(jobChan)
	}()
	go func() {
		consumerWG.Wait()
		close(resultChan)
	}()
	total := producerCount * itemsPerProducer
	results := make([]int, 0, total)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result, ok := <-resultChan:
			if !ok {
				return results, nil
			}
			results = append(results, result)
		}
	}
}
