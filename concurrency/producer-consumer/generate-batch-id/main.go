package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	results, err := ProcessBatchIDs(
		ctx,
		2, // producerCount
		3, // consumerCount
		4, // itemsPerProducer
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(
		"Expected values:",
		[]int{1050, 1150, 1250, 1350, 1450, 1550, 1650, 1750},
	)
	fmt.Println("Actual values:  ", results)
}
