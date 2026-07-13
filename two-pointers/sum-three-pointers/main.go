package main

import (
	"fmt"
	"sort"
)

func ThreeSum(nums []int) [][]int {
	// Em implement phần này.
	result := make([][]int, 0)
	sort.Slice(nums, func(i, j int) bool { return nums[i] < nums[j] })
	for i := 0; i <= len(nums)-2; i++ {
		left := i + 1
		right := len(nums) - 1
		sum := nums[i] + nums[left] + nums[right]
		for left < right {
			if nums[i] > 0 {
				break
			}
			if sum < 0 {
				left++
			}
			if sum > 0 {
				right--
			}
			if sum == 0 {
				result = append(result, []int{nums[i], nums[left], nums[right]})
			}
			left++
			right--
			for left < right && nums[left] == nums[left-1] {
				left++
			}
			for left < right && nums[right] == nums[right-1] {
				right--
			}

		}
	}

	return result
}

func main() {
	testCases := []struct {
		nums     []int
		expected [][]int
	}{
		{
			nums:     []int{-1, 0, 1, 2, -1, -4},
			expected: [][]int{{-1, -1, 2}, {-1, 0, 1}},
		},
		{
			nums:     []int{0, 1, 1},
			expected: [][]int{},
		},
		{
			nums:     []int{0, 0, 0},
			expected: [][]int{{0, 0, 0}},
		},
	}

	for _, tc := range testCases {
		actual := ThreeSum(tc.nums)

		fmt.Printf(
			"nums=%v\nexpected=%v\nactual=%v\n\n",
			tc.nums,
			tc.expected,
			actual,
		)
	}
}
