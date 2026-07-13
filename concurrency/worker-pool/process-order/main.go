package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	orders := []int{100, 200, 300, 400}

	results, err := ProcessOrders(ctx, orders, 2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Expected:", []int{110, 220, 330, 440})
	fmt.Println("Actual:  ", results)
}
