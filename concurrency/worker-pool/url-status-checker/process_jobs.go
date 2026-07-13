package main

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"
)

// Problem: Concurrent URL Status Checker
//
// You are given a list of URLs and a fixed number of workers.
//
// Each worker receives a URL from a shared job channel,
// performs an HTTP GET request, and sends the result to a shared result channel.
//
// Requirements:
//
//  1. Create exactly workerCount worker goroutines.
//  2. Send URLs through a channel.
//  3. Each worker performs an HTTP GET request.
//  4. Return one result for each URL.
//  5. Preserve the original order of the input URLs.
//  6. Each result must contain:
//     - URL
//     - HTTP status code
//     - error, if the request failed
//  7. Support context cancellation.
//  8. Return an error when workerCount <= 0.
//  9. Do not use time.Sleep for synchronization.
//
// 10. Do not leak goroutines.
//
// Function:
//
//	func CheckURLs(
//		ctx context.Context,
//		urls []string,
//		workerCount int,
//	) ([]URLResult, error)
//
// Types:
//
//	type URLResult struct {
//		URL        string
//		StatusCode int
//		Err        error
//	}
//
// Example:
//
//	Input:
//		urls = []string{
//			"https://example.com",
//			"https://httpstat.us/404",
//		}
//		workerCount = 2
//
//	Expected:
//		[]URLResult{
//			{URL: "https://example.com", StatusCode: 200},
//			{URL: "https://httpstat.us/404", StatusCode: 404},
//		}
//
// Notes:
//
//	The result order must match the input URL order,
//	even if requests finish in a different order.
type URLResult struct {
	Index      int
	URL        string
	StatusCode int
	Err        error
}
type URLJob struct {
	Index int
	URL   string
}

func CheckURLs(
	ctx context.Context,
	urls []string,
	workerCount int,
) ([]URLResult, error) {
	// TODO: implement worker pool
	if len(urls) == 0 {
		return []URLResult{}, nil
	}
	if workerCount <= 0 {
		return nil, errors.New("workerCount must be greater than zero")
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resultChan := make(chan URLResult, len(urls))
	jobChan := make(chan URLJob, len(urls))
	var wg sync.WaitGroup
	worker := func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case job, ok := <-jobChan:
				if !ok {
					return
				}
				result := callUrl(ctx, client, job)
				select {
				case <-ctx.Done():
					return
				case resultChan <- result:
				}
			}
		}
	}
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker()
	}
	go func() {
		defer close(jobChan)
		for index, url := range urls {
			select {
			case <-ctx.Done():
				return
			case jobChan <- URLJob{URL: url, Index: index}:
			}
		}
	}()
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	results := make([]URLResult, len(urls))
	recived := 0
	for recived < len(urls) {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		case result, ok := <-resultChan:
			if !ok {
				return results, nil
			}
			results[result.Index] = result
			recived++
		}
	}
	return results, nil
}
func callUrl(ctx context.Context, client *http.Client, job URLJob) URLResult {
	result := URLResult{
		Index: job.Index,
		URL:   job.URL,
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, job.URL, nil)
	if err != nil {
		result.Err = err
		return result
	}
	resp, err := client.Do(req)
	if err != nil {
		result.Err = err
		return result
	}
	defer resp.Body.Close()
	result.StatusCode = resp.StatusCode
	return result
}
