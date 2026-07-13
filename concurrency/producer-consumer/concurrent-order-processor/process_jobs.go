package main

import (
	"context"
	"fmt"
	"sync"
)

// Problem: Concurrent Order Processor
//
// You are given multiple producers and multiple consumers.
//
// Each producer generates order amounts and sends them into a shared order channel.
// Each consumer receives an order amount, applies a 10% service fee,
// and sends the final amount into a shared result channel.
//
// Requirements:
//
//  1. Create exactly producerCount producer goroutines.
//
//  2. Create exactly consumerCount consumer goroutines.
//
//  3. Each producer generates exactly ordersPerProducer orders.
//
//  4. Order values must be unique and start from 100.
//
//  5. Each consumer calculates:
//
//     finalAmount = orderAmount + orderAmount/10
//
//  6. Close orderCh only after all producers finish.
//
//  7. Close resultCh only after all consumers finish.
//
//  8. Support context cancellation.
//
//  9. Do not use time.Sleep.
//
// 10. Do not leak goroutines.
//
// Function:
//
//	func ProcessOrders(
//		ctx context.Context,
//		producerCount int,
//		consumerCount int,
//		ordersPerProducer int,
//	) ([]int, error)
//
// Example:
//
//	Input:
//		producerCount    = 2
//		consumerCount    = 3
//		ordersPerProducer = 3
//
// Producers generate:
//
//	Producer 0: 100, 110, 120
//	Producer 1: 130, 140, 150
//
// Expected values:
//
//	110, 121, 132, 143, 154, 165
//
// Notes:
//
//	The result order is not guaranteed.
func ProcessOrders(
	ctx context.Context,
	producerCount int,
	consumerCount int,
	ordersPerProducer int,
) ([]int, error) {
	// TODO
	if producerCount <= 0 || consumerCount <= 0 {
		return nil, fmt.Errorf("producer count and consumer count must be positive")
	}
	if ordersPerProducer <= 0 {
		return nil, fmt.Errorf("producer count must be positive")
	}
	reslutChan := make(chan int)
	jobChan := make(chan int)
	var producerWG sync.WaitGroup
	var consumerWG sync.WaitGroup
	producer := func(producerID int) {
		defer producerWG.Done()
		for i := 0; i < ordersPerProducer; i++ {
			globalIndex := producerID*ordersPerProducer + i
			orderAmount := 100 + globalIndex*10
			select {
			case <-ctx.Done():
				return
			case jobChan <- orderAmount:
			}
		}
	}
	consumer := func() {
		defer consumerWG.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case orderAmount, ok := <-jobChan:
				if !ok {
					return
				}
				finalAmount := orderAmount + orderAmount/10
				select {
				case <-ctx.Done():
					return
				case reslutChan <- finalAmount:
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
		close(reslutChan)
	}()
	total := producerCount * ordersPerProducer
	results := make([]int, 0, total)
	for {
		select {
		case <-ctx.Done():
			return results, nil
		case result, ok := <-reslutChan:
			if !ok {
				return results, nil
			}
			results = append(results, result)
		}
	}
}
