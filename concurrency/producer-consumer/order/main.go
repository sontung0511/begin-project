package main

import (
	"context"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Starting producer")
	ctx := context.Background()

	orders := []int{100, 200, 300, 400}
	workerCount := 2

	result, err := ProcessOrders(ctx, orders, workerCount)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
