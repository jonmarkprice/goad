package main

import (
	"testing"

	// "github.com/goadapp/goad/lambda"
)

func TestBin(t *testing.T) {
	counts := make(map[int64]int)

	if len(counts) != 0 {
		t.Error("Counts should be empty initially")
	}

	Bin(11, counts, 5)

	if counts[10] != 1 {
		t.Error("11 should be binned into 10 (by 5s).")
	}

	if len(counts) != 1 {
		t.Error("Counts should now be 1.")
	}
}
