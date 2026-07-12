package main

import (
	"fmt"
	"sync"
)

// PassMessage sends an initial message through a sequence of concurrency.
//
// Each station runs in its own goroutine and performs three steps:
//  1. Receives the message from the previous station.
//  2. Appends its station number to the message.
//  3. Sends the updated message to the next station.
//
// Example:
//
//	Input:  concurrency = 5
//	Output: Start -> Station 1 -> Station 2 -> Station 3 -> Station 4 -> Station 5
//
// Implementation notes:
//   - Use one goroutine for each station.
//   - Use channels to pass the message between concurrency.
//   - For n concurrency, create n+1 channels.
//   - Station i receives from channels[i-1].
//   - Station i sends to channels[i].
//   - Use sync.WaitGroup to wait for all goroutines to finish.
//   - Return an empty string when concurrency <= 0.

func PassMessage(stations int) string {
	// TODO:
	// 1. Nếu concurrency <= 0 thì return ""
	if stations <= 0 {
		return ""
	}
	// 2. Tạo concurrency + 1 channel
	channels := make([]chan string, stations+1)
	// 3. Khởi tạo từng channel
	for i := 0; i <= stations; i++ {
		channels[i] = make(chan string)
	}
	var wg sync.WaitGroup
	for i := 1; i <= stations; i++ {
		wg.Add(1)
		// 4. Tạo 1 goroutine cho mỗi station
		go func(station int) {
			defer wg.Done()
			// 5. Station i nhận từ channels[i-1]
			message := <-channels[station-1]
			// 6. Thêm nội dung: " -> Station i"
			message = fmt.Sprintf("%s -> Station %d", message, station)
			// 7. Gửi sang channels[i]
			channels[station] <- message
		}(i)
	}
	go func() {
		// 8. Main gửi "Start" vào channels[0]
		channels[0] <- "Start"
	}()
	// 9. Main nhận kết quả từ channels[concurrency]
	finalMessage := <-channels[stations]
	// 10. Chờ tất cả goroutine hoàn thành
	wg.Wait()

	fmt.Println("All concurrency had receive Final Message:", finalMessage)
	return finalMessage
}
