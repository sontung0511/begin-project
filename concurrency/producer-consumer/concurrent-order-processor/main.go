package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	results, err := ProcessOrders(
		ctx,
		2, // producerCount
		3, // consumerCount
		3, // ordersPerProducer
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Expected values:", []int{110, 121, 132, 143, 154, 165})
	fmt.Println("Actual values:  ", results)
}
