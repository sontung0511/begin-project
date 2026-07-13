package main

import "fmt"

func SumIndex(nums []int, target int) []int {
	mapNums := make(map[int]int)
	for i, num := range nums {
		needNum := target - num
		if idx, ok := mapNums[needNum]; ok {
			return []int{idx, i}
		}
		mapNums[num] = i
	}
	return nil
}
func main() {
	nums := []int{1, 2, 3, 4, 5}
	target := 5
	result := SumIndex(nums, target)
	fmt.Println(result)
}
