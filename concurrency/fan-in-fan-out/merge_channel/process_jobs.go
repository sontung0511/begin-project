package main

import (
	"context"
	"sync"
)

// Problem: Merge Multiple Channels
//
// You are given multiple read-only integer channels.
//
// Create a function that reads values from all input channels concurrently
// and merges them into one output channel.
//
// Requirements:
//
//  1. Start one goroutine for each input channel.
//  2. Forward every value from each input channel into one shared output channel.
//  3. Close the output channel only after all input channels are fully consumed.
//  4. Do not close any input channel.
//  5. Do not use time.Sleep for synchronization.
//  6. Do not leak goroutines.
//  7. Support context cancellation.
//  8. When ctx is cancelled, stop all forwarding goroutines and close output.
//
// Function:
//
//	func MergeChannels(
//		ctx context.Context,
//		channels ...<-chan int,
//	) <-chan int
//
// Example:
//
//	ch1 := make(chan int)
//	ch2 := make(chan int)
//
//	go func() {
//		defer close(ch1)
//		ch1 <- 1
//		ch1 <- 3
//		ch1 <- 5
//	}()
//
//	go func() {
//		defer close(ch2)
//		ch2 <- 2
//		ch2 <- 4
//		ch2 <- 6
//	}()
//
//	output := MergeChannels(context.Background(), ch1, ch2)
//
//	for value := range output {
//		fmt.Println(value)
//	}
//
// Expected:
//
//	All values 1, 2, 3, 4, 5, 6 must appear exactly once.
//
// Notes:
//
//	The output order is not guaranteed because the input channels are read
//	concurrently.

func MergeChannels(
	ctx context.Context,
	channels ...<-chan int,
) <-chan int {
	// TODO: implement fan-in
	var wg sync.WaitGroup
	out := make(chan int)
	wg.Add(len(channels))
	for _, inputCh := range channels {
		go func(ch <-chan int) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case value, ok := <-ch:
					if !ok {
						return
					}

					select {
					case <-ctx.Done():
						return
					case out <- value:
					}
				}
			}
		}(inputCh)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
