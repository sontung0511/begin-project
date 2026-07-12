package main

import "testing"

func TestExample(t *testing.T) {
	got := 1 + 1
	want := 2

	if got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}
