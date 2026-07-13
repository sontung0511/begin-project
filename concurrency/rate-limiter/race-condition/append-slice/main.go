package main

//
//import (
//	"fmt"
//	"sync"
//)
//
//func main() {
//	var wg sync.WaitGroup
//	var mu sync.Mutex
//	results := make([]int, 0)
//
//	for i := 0; i < 1000; i++ {
//		wg.Add(1)
//
//		go func(value int) {
//			defer wg.Done()
//			mu.Lock()
//			results = append(results, value)
//			mu.Unlock()
//		}(i)
//	}
//
//	wg.Wait()
//
//	fmt.Println("Expected length:", 1000)
//	fmt.Println("Actual length:  ", len(results))
//}
