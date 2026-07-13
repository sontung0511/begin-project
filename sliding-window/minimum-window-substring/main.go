package main

import "fmt"

func minWindow(s string, t string) string {
	if len(t) == 0 || len(s) < len(t) {
		return ""
	}

	need := make(map[byte]int)
	for i := 0; i < len(t); i++ {
		need[t[i]]++
	}

	window := make(map[byte]int)

	required := len(need)
	formed := 0

	left := 0
	minLen := len(s) + 1
	start := 0

	for right := 0; right < len(s); right++ {
		char := s[right]
		window[char]++

		if need[char] > 0 && window[char] == need[char] {
			formed++
		}

		for formed == required {
			currentLen := right - left + 1

			if currentLen < minLen {
				minLen = currentLen
				start = left
			}

			leftChar := s[left]
			window[leftChar]--

			if need[leftChar] > 0 &&
				window[leftChar] < need[leftChar] {
				formed--
			}

			left++
		}
	}

	if minLen == len(s)+1 {
		return ""
	}

	return s[start : start+minLen]
}

func main() {
	s := "OUZODYXAZV"
	t := "XYZ"

	result := minWindow(s, t)

	fmt.Println(result) // YXAZ
}
