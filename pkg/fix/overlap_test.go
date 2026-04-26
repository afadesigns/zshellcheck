// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package fix

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
)

func TestOverlapSameLineOverlapping(t *testing.T) {
	a := katas.FixEdit{Line: 1, Column: 1, Length: 5, Replace: ""}
	b := katas.FixEdit{Line: 1, Column: 3, Length: 4, Replace: ""}
	if !Overlap(a, b) {
		t.Errorf("expected overlap on same line")
	}
}

func TestOverlapSameLineDisjoint(t *testing.T) {
	a := katas.FixEdit{Line: 1, Column: 1, Length: 2, Replace: ""}
	b := katas.FixEdit{Line: 1, Column: 5, Length: 2, Replace: ""}
	if Overlap(a, b) {
		t.Errorf("disjoint edits reported as overlapping")
	}
}

func TestOverlapDifferentLines(t *testing.T) {
	a := katas.FixEdit{Line: 1, Column: 1, Length: 5, Replace: ""}
	b := katas.FixEdit{Line: 2, Column: 1, Length: 5, Replace: ""}
	if Overlap(a, b) {
		t.Errorf("edits on different lines reported as overlapping")
	}
}

func TestDiff_EmptyEditsReturnsBlank(t *testing.T) {
	out, err := Diff("file.zsh", "echo hi\n", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty diff for no edits, got %q", out)
	}
}

func TestDiff_NonEmptyChange(t *testing.T) {
	src := "echo hi\n"
	edits := []katas.FixEdit{{Line: 1, Column: 1, Length: 4, Replace: "print"}}
	out, err := Diff("file.zsh", src, edits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Errorf("expected non-empty diff")
	}
}

func TestDiff_PropagatesApplyError(t *testing.T) {
	edits := []katas.FixEdit{{Line: 99, Column: 1, Length: 1, Replace: ""}}
	if _, err := Diff("file.zsh", "echo hi\n", edits); err == nil {
		t.Errorf("expected error for out-of-range edit")
	}
}

func TestApply_NegativeLengthRejected(t *testing.T) {
	edits := []katas.FixEdit{{Line: 1, Column: 1, Length: -1, Replace: ""}}
	if _, err := Apply("echo hi\n", edits); err == nil {
		t.Errorf("expected error for negative length")
	}
}

func TestApply_NonPositiveColumnRejected(t *testing.T) {
	edits := []katas.FixEdit{{Line: 1, Column: 0, Length: 1, Replace: ""}}
	if _, err := Apply("echo hi\n", edits); err == nil {
		t.Errorf("expected error for non-positive column")
	}
}

func TestApply_StartPastEndRejected(t *testing.T) {
	edits := []katas.FixEdit{{Line: 1, Column: 999, Length: 0, Replace: ""}}
	if _, err := Apply("echo hi\n", edits); err == nil {
		t.Errorf("expected error for start past end of source")
	}
}

func TestApply_EndPastEndRejected(t *testing.T) {
	edits := []katas.FixEdit{{Line: 1, Column: 1, Length: 999, Replace: ""}}
	if _, err := Apply("echo hi\n", edits); err == nil {
		t.Errorf("expected error for end past end of source")
	}
}
