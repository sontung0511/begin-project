package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	results, err := GenerateOddNumbers(
		ctx,
		2, // producerCount
		3, // consumerCount
		3, // itemsPerProducer
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Expected values:", []int{1, 3, 5, 7, 9, 11})
	fmt.Println("Actual values:  ", results)
}
