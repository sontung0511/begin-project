package main

import (
	"context"
	"fmt"
	"sync"
)

// Problem: Producer–Consumer
//
// You are given multiple producers and multiple consumers.
//
// Each producer generates integer jobs and sends them into one shared job channel.
// Each consumer reads jobs from that channel, processes them, and sends results
// into one shared result channel.
//
// Requirements:
//
//  1. Create exactly producerCount producer goroutines.
//  2. Create exactly consumerCount consumer goroutines.
//  3. Producers must send jobs through a shared channel.
//  4. Consumers must read jobs from that shared channel.
//  5. Each consumer calculates job * job.
//  6. Close the job channel only after all producers finish.
//  7. Close the result channel only after all consumers finish.
//  8. Support context cancellation.
//  9. Do not use time.Sleep for synchronization.
//
// 10. Do not leak goroutines.
//
// Function:
//
//	func ProducerConsumer(
//		ctx context.Context,
//		producerCount int,
//		consumerCount int,
//		jobsPerProducer int,
//	) ([]int, error)
//
// Example:
//
//	Input:
//		producerCount  = 2
//		consumerCount  = 3
//		jobsPerProducer = 3
//
// Producers generate:
//
//	Producer 0: 1, 2, 3
//	Producer 1: 4, 5, 6
//
// Expected results:
//
//	1, 4, 9, 16, 25, 36
//
// Notes:
//
//	The result order is not guaranteed because consumers process concurrently.
func ProducerConsumer(
	ctx context.Context,
	producerCount int,
	consumerCount int,
	jobsPerProducer int,
) ([]int, error) {
	if producerCount < 1 {
		return nil, fmt.Errorf("producer count must be greater than zero")
	}
	if consumerCount < 1 {
		return nil, fmt.Errorf("consumer count must be greater than zero")
	}
	if jobsPerProducer <= 0 {
		return nil, fmt.Errorf("jobs per producer must be greater than zero")
	}
	resultCh := make(chan int, producerCount)
	jobCh := make(chan int, consumerCount)
	var producerWG sync.WaitGroup
	var consumerWG sync.WaitGroup
	// TODO 1: Producer function
	producer := func(producerID int) {
		defer producerWG.Done()
		// Mỗi producer tạo jobsPerProducer job.
		// Công thức gợi ý:
		// start := producerID*jobsPerProducer + 1
		start := producerID*jobsPerProducer + 1
		// TODO:
		// - Duyệt jobsPerProducer lần
		for i := 0; i < jobsPerProducer; i++ {
			// - Tạo job
			job := start + i
			select {
			case <-ctx.Done():
				return
			case jobCh <- job:
			}
		}
		// - select giữa ctx.Done() và jobCh <- job
	}
	// TODO 2: Consumer function
	consumer := func() {
		defer consumerWG.Done()

		// TODO:
		// - Liên tục nhận job từ jobCh
		for {
			// - Nếu context cancel thì return
			select {
			case <-ctx.Done():
				return
				// - Nếu jobCh đóng thì return
			case job, ok := <-jobCh:
				if !ok {
					return
				}
				// - Tính square := job * job
				// - Gửi square vào resultCh
				square := job * job
				// - Khi gửi cũng phải nghe ctx.Done()
				select {
				case <-ctx.Done():
					return
				case resultCh <- square:
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
		close(jobCh)
	}()
	go func() {
		consumerWG.Wait()
		close(resultCh)
	}()
	total := producerCount * jobsPerProducer
	results := make([]int, 0, total)
	for {
		select {
		case <-ctx.Done():
			return results, nil
		case result, ok := <-resultCh:
			if !ok {
				return results, nil
			}
			results = append(results, result)
		}
	}
}
