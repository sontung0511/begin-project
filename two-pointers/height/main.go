package main

import "fmt"

func MaxArea(height []int) int {
	// Em implement phần này.
	if len(height) == 0 {
		return 0
	}
	left, right := 0, len(height)-1
	maxArea := 0
	for left < right {
		h := height[left]
		if h > height[right] {
			h = height[right]
		}
		width := right - left
		area := h * width
		if area > maxArea {
			maxArea = area
		}
		if height[left] < height[right] {
			left++
		} else {
			right--
		}
	}
	return maxArea
}

func main() {
	testCases := []struct {
		height   []int
		expected int
	}{
		{
			height:   []int{1, 8, 6, 2, 5, 4, 8, 3, 7},
			expected: 49,
		},
		{
			height:   []int{1, 1},
			expected: 1,
		},
		{
			height:   []int{1, 7, 2, 5, 4, 7, 3, 6},
			expected: 36,
		},
	}

	for _, tc := range testCases {
		actual := MaxArea(tc.height)

		fmt.Printf(
			"height=%v expected=%d actual=%d passed=%v\n",
			tc.height,
			tc.expected,
			actual,
			tc.expected == actual,
		)
	}
}
