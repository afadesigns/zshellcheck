// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/token"
)

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

// TestLongFlagName pins the split between a long option and its inline
// value, including the words that are not long options at all.
func TestLongFlagName(t *testing.T) {
	cases := []struct {
		word      string
		wantName  string
		wantValue bool
	}{
		{"--depth=1", "--depth", true},
		{"--depth", "--depth", false},
		{"--filter=blob:none", "--filter", true},
		{`--depth="1"`, "--depth", true},
		{"--depth=", "--depth", true},
		{"--", "--", false},
		{"-d=1", "-d=1", false},
		{"clone", "clone", false},
		{"", "", false},
	}
	for _, tc := range cases {
		name, value := LongFlagName(tc.word)
		if name != tc.wantName || value != tc.wantValue {
			t.Errorf("LongFlagName(%q) = (%q, %v), want (%q, %v)",
				tc.word, name, value, tc.wantName, tc.wantValue)
		}
	}
}

// TestZC1231LimitsHistory pins which `git clone` options count as already
// limiting the downloaded history.
func TestZC1231LimitsHistory(t *testing.T) {
	for _, word := range []string{
		"--depth", "--depth=1", "--shallow-since=2024-01-01",
		"--shallow-since", "--shallow-exclude=refs/tags/v1",
	} {
		if !zc1231LimitsHistory(word) {
			t.Errorf("zc1231LimitsHistory(%q) = false, want true", word)
		}
	}
	for _, word := range []string{
		"--single-branch", "--filter=blob:none", "--recursive", "-b", "clone", "",
	} {
		if zc1231LimitsHistory(word) {
			t.Errorf("zc1231LimitsHistory(%q) = true, want false", word)
		}
	}
}

// TestFixZC1231RefusesSecondShallowFlag pins the guard inside the fix. The
// check never reports a clone that already limits its history, so this path
// exists to keep a future change to the check from producing
// `--depth 1 --shallow-since=...`, which git refuses outright.
func TestFixZC1231RefusesSecondShallowFlag(t *testing.T) {
	// The fix locates the `clone` token in the source by position, so the
	// argument carries the column it occupies in `source`.
	clone := func(flag string) *ast.SimpleCommand {
		cloneArg := &ast.Identifier{
			Token: token.Token{Literal: "clone", Line: 1, Column: 5},
			Value: "clone",
		}
		args := []ast.Expression{cloneArg}
		if flag != "" {
			args = append(args, &ast.Identifier{Value: flag})
		}
		args = append(args, &ast.Identifier{Value: "https://example.com/r"})
		return &ast.SimpleCommand{Name: &ast.Identifier{Value: "git"}, Arguments: args}
	}
	source := []byte("git clone https://example.com/r\n")
	for _, flag := range []string{"--depth=1", "--shallow-since=2024-01-01", "--shallow-exclude=v1"} {
		if edits := fixZC1231(clone(flag), Violation{}, source); edits != nil {
			t.Errorf("fixZC1231 with %s returned %d edits, want none", flag, len(edits))
		}
	}
	// A clone with no such flag still gets the insertion.
	if edits := fixZC1231(clone(""), Violation{}, source); len(edits) != 1 {
		t.Errorf("fixZC1231 on a full clone returned %d edits, want 1", len(edits))
	}
}
