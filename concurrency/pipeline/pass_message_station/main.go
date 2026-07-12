package main

import (
	"fmt"
)

func main() {
	result := PassMessage(5)

	fmt.Println("Expected:")
	fmt.Println("Start -> Station 1 -> Station 2 -> Station 3 -> Station 4 -> Station 5")

	fmt.Println("Actual:")
	fmt.Println(result)
}
