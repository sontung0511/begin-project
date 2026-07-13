package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	urls := []string{
		"https://example.com",
		"https://httpstat.us/404",
		"https://httpstat.us/500",
	}

	results, err := CheckURLs(ctx, urls, 3)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, result := range results {
		fmt.Printf(
			"URL: %s | Status: %d | Error: %v\n",
			result.URL,
			result.StatusCode,
			result.Err,
		)
	}
}
