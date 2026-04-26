// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import "testing"

func TestLineColToByteOffset(t *testing.T) {
	src := []byte("abc\ndef\nghi")
	cases := []struct {
		line, col int
		want      int
	}{
		{1, 1, 0},  // 'a'
		{1, 3, 2},  // 'c'
		{1, 4, 3},  // newline after 'c'
		{2, 1, 4},  // 'd'
		{3, 3, 10}, // 'i'
		{3, 4, 11}, // EOF on last line (no trailing newline)
	}
	for _, c := range cases {
		got := LineColToByteOffset(src, c.line, c.col)
		if got != c.want {
			t.Errorf("LineColToByteOffset(%d,%d)=%d, want %d", c.line, c.col, got, c.want)
		}
	}
}

func TestLineColToByteOffset_OutOfRange(t *testing.T) {
	src := []byte("one")
	if LineColToByteOffset(src, 99, 1) != -1 {
		t.Error("line past end should return -1")
	}
	if LineColToByteOffset(src, 0, 1) != -1 {
		t.Error("line 0 should return -1")
	}
	if LineColToByteOffset(src, 1, 0) != -1 {
		t.Error("col 0 should return -1")
	}
}

func TestIdentLenAt(t *testing.T) {
	src := []byte("which foo; print bar")
	if got := IdentLenAt(src, 0); got != 5 {
		t.Errorf("IdentLenAt(0)=%d, want 5 (\"which\")", got)
	}
	if got := IdentLenAt(src, 6); got != 3 {
		t.Errorf("IdentLenAt(6)=%d, want 3 (\"foo\")", got)
	}
	if got := IdentLenAt(src, 5); got != 0 {
		t.Errorf("IdentLenAt(5)=%d, want 0 (space)", got)
	}
	if got := IdentLenAt(src, len(src)); got != 0 {
		t.Errorf("IdentLenAt(eof)=%d, want 0", got)
	}
}
