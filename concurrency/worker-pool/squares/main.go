package main

import (
	"context"
	"fmt"
	"log"
)

func main() {
	ctx := context.Background()

	jobs := []int{1, 2, 3, 4, 5}

	results, err := ProcessJobs(ctx, jobs, 3)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Expected:")
	fmt.Println([]int{1, 4, 9, 16, 25})

	fmt.Println("Actual:")
	fmt.Println(results)
}
