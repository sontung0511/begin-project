package main

import "fmt"

func trap(height []int) int {
	left, right := 0, len(height)-1
	totalArea := 0
	maxRight := 0
	maxLeft := 0
	for left < right {
		if height[left] < height[right] {
			if height[left] > maxLeft {
				maxLeft = height[left]
			} else {
				totalArea += maxLeft - height[left]
			}
			left++
		} else {
			if height[right] > maxRight {
				maxRight = height[right]
			} else {
				totalArea += maxRight - height[right]
			}
			right--
		}
	}
	return totalArea
}
func main() {
	height := []int{0, 2, 0, 3, 1, 0, 1, 3, 2, 1}
	area := trap(height)
	fmt.Println(area)
}
