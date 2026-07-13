package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	numbers := []int{1, 2, 3, 4, 5}
	workerCount := 3

	result, err := SumSquares(ctx, numbers, workerCount)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Expected:", 55)
	fmt.Println("Actual:  ", result)
}
