package main

import (
	"fmt"
	"sync"
)

type User struct {
	Name string
	Age  int
}

func main() {
	var wg sync.WaitGroup
	var mu sync.RWMutex
	user := User{
		Name: "Tyler",
		Age:  20,
	}

	wg.Add(2)

	go func() {
		defer wg.Done()
		mu.Lock()
		for i := 0; i < 1000; i++ {
			user.Age++
		}
		mu.Unlock()

	}()

	go func() {
		defer wg.Done()

		for i := 0; i < 1000; i++ {
			mu.RLock()
			age := user.Age
			mu.RUnlock()
			fmt.Println(age)
		}
	}()

	wg.Wait()
}
