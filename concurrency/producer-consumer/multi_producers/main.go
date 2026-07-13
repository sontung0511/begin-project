package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	results, err := ProducerConsumer(ctx, 2, 3, 3)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Results:", results)
}
