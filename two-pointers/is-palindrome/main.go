package main

import (
	"fmt"
	"strings"
	"unicode"
)

func IsPalindrome(s string) bool {
	s = strings.ToLower(s)
	// Em implement phần này.
	left, right := 0, len(s)-1
	for left < right {
		for left < right && !unicode.IsLetter(rune(s[left])) && !unicode.IsDigit(rune(s[left])) {
			left++
		}
		for left < right && !unicode.IsLetter(rune(s[right])) && !unicode.IsDigit(rune(s[right])) {
			right--
		}
		if s[left] != s[right] {
			return false
		}
		left++
		right--
	}
	return true
}

func main() {
	testCases := []struct {
		input    string
		expected bool
	}{
		{
			input:    "A man, a plan, a canal: Panama",
			expected: true,
		},
		{
			input:    "race a car",
			expected: false,
		},
		{
			input:    " ",
			expected: true,
		},
		{
			input:    "0P",
			expected: false,
		},
	}

	for _, tc := range testCases {
		actual := IsPalindrome(tc.input)

		fmt.Printf(
			"input=%q expected=%v actual=%v passed=%v\n",
			tc.input,
			tc.expected,
			actual,
			tc.expected == actual,
		)
	}
}
