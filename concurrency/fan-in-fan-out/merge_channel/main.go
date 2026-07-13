package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		defer close(ch1)

		ch1 <- 1
		ch1 <- 3
		ch1 <- 5
	}()

	go func() {
		defer close(ch2)

		ch2 <- 2
		ch2 <- 4
		ch2 <- 6
	}()

	output := MergeChannels(ctx, ch1, ch2)

	for value := range output {
		fmt.Println(value)
	}
}
