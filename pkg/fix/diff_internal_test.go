// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package fix

import (
	"strings"
	"testing"
)

func TestSplitLines_Empty(t *testing.T) {
	if got := splitLines(""); got != nil {
		t.Errorf("splitLines(\"\") = %v, want nil", got)
	}
}

func TestSplitLines_TrailingNewline(t *testing.T) {
	got := splitLines("a\nb\n")
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Errorf("unexpected: %v", got)
	}
}

func TestSplitLines_NoTrailingNewline(t *testing.T) {
	got := splitLines("a\nb")
	if len(got) != 2 {
		t.Errorf("expected 2 lines, got %d", len(got))
	}
}

func TestMaxIntMinInt(t *testing.T) {
	if maxInt(1, 2) != 2 || maxInt(5, 3) != 5 {
		t.Error("maxInt wrong")
	}
	if minInt(1, 2) != 1 || minInt(5, 3) != 3 {
		t.Error("minInt wrong")
	}
}

func TestUnifiedDiff_Insertion(t *testing.T) {
	out := unifiedDiff("file", "a\nb\n", "a\nx\nb\n")
	if !strings.Contains(out, "+x") {
		t.Errorf("expected insertion line, got %q", out)
	}
}

func TestUnifiedDiff_Deletion(t *testing.T) {
	out := unifiedDiff("file", "a\nb\nc\n", "a\nc\n")
	if !strings.Contains(out, "-b") {
		t.Errorf("expected deletion line, got %q", out)
	}
}

func TestUnifiedDiff_Replacement(t *testing.T) {
	out := unifiedDiff("file", "a\nold\nc\n", "a\nnew\nc\n")
	if !strings.Contains(out, "-old") || !strings.Contains(out, "+new") {
		t.Errorf("expected replacement, got %q", out)
	}
}

func TestUnifiedDiff_TailEdit(t *testing.T) {
	out := unifiedDiff("file", "a\nb\nc\n", "a\nb\nc\nd\n")
	if !strings.Contains(out, "+d") {
		t.Errorf("expected tail insertion, got %q", out)
	}
}

func TestLcsTable_Square(t *testing.T) {
	a := []string{"x", "y"}
	b := []string{"a", "b"}
	tbl := lcsTable(a, b)
	if tbl[0][0] != 0 {
		t.Errorf("expected 0 lcs for disjoint slices, got %d", tbl[0][0])
	}
}
