package main

import "fmt"

func TwoSumPointer(numbers []int, target int) []int {
	left, right := 0, len(numbers)-1
	for left < right {
		sum := numbers[left] + numbers[right]
		if sum == target {
			return []int{left + 1, right + 1}
		}
		if sum < target {
			left++
		}
		if sum > target {
			right--
		}
	}
	return nil
}
func main() {
	nums := []int{2, 4, 5, 6, 7}
	target := 6
	result := TwoSumPointer(nums, target)
	fmt.Println(result)
}
