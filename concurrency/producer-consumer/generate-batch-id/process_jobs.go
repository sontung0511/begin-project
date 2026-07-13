package main

import (
	"context"
	"fmt"
	"sync"
)

// Problem:
//
// Multiple producers generate unique batch IDs.
// Each consumer receives a batch ID and adds a processing fee of 50.
//
// Formula:
//
// globalIndex := producerID*itemsPerProducer + i
// batchID := 1000 + globalIndex*100
// processedValue := batchID + 50
//
// Example:
//
// producerCount = 2
// consumerCount = 3
// itemsPerProducer = 4
//
// Generated batch IDs:
//
// Producer 0: 1000, 1100, 1200, 1300
// Producer 1: 1400, 1500, 1600, 1700
//
// Expected processed values:
//
// 1050, 1150, 1250, 1350, 1450, 1550, 1650, 1750
//
// Result order is not guaranteed.

func ProcessBatchIDs(
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
			order := 1000 + start*100
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
				select {
				case <-ctx.Done():
					return
				case resultChan <- job + 50:
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
	results := make([]int, 0, producerCount*itemsPerProducer)
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
