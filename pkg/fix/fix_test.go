package fix

import (
	"strings"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
)

func TestApply_Empty(t *testing.T) {
	got, err := Apply("hello\n", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "hello\n" {
		t.Fatalf("expected passthrough, got %q", got)
	}
}

func TestApply_SingleEdit(t *testing.T) {
	src := "echo hi\n"
	edits := []katas.FixEdit{{Line: 1, Column: 1, Length: 4, Replace: "print"}}
	got, err := Apply(src, edits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "print hi\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestApply_MultipleEditsOnSameLine(t *testing.T) {
	src := "foo bar baz\n"
	// Replace "foo" with "XX" and "baz" with "YYYY". Order should not matter.
	edits := []katas.FixEdit{
		{Line: 1, Column: 1, Length: 3, Replace: "XX"},
		{Line: 1, Column: 9, Length: 3, Replace: "YYYY"},
	}
	got, err := Apply(src, edits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "XX bar YYYY\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestApply_EditsAcrossLines(t *testing.T) {
	src := "one\ntwo\nthree\n"
	edits := []katas.FixEdit{
		{Line: 1, Column: 1, Length: 3, Replace: "ONE"},
		{Line: 3, Column: 1, Length: 5, Replace: "THREE"},
	}
	got, err := Apply(src, edits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "ONE\ntwo\nTHREE\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestApply_Deletion(t *testing.T) {
	src := "keep delete keep\n"
	edits := []katas.FixEdit{{Line: 1, Column: 6, Length: 7, Replace: ""}}
	got, err := Apply(src, edits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "keep keep\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestApply_MultiLineReplacement(t *testing.T) {
	src := "before\nmiddle\nafter\n"
	edits := []katas.FixEdit{{Line: 2, Column: 1, Length: 6, Replace: "a\nb\nc"}}
	got, err := Apply(src, edits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "before\na\nb\nc\nafter\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestApply_OverlappingEdits(t *testing.T) {
	// The outer span ("hello") should win; the inner (overlap) edit
	// is dropped silently. Running Apply again on the output would
	// then be a no-op because the new text no longer contains the
	// inner pattern.
	src := "hello world\n"
	edits := []katas.FixEdit{
		{Line: 1, Column: 1, Length: 5, Replace: "HI"},
		{Line: 1, Column: 4, Length: 3, Replace: "XX"},
	}
	got, err := Apply(src, edits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "HI world\n" {
		t.Fatalf("got %q, want %q", got, "HI world\n")
	}
}

func TestApply_OutOfRangeLine(t *testing.T) {
	_, err := Apply("one\n", []katas.FixEdit{{Line: 99, Column: 1, Length: 1, Replace: "x"}})
	if err == nil {
		t.Fatal("expected out-of-range error")
	}
}

func TestApply_OutOfRangeColumn(t *testing.T) {
	_, err := Apply("abc\n", []katas.FixEdit{{Line: 1, Column: 1, Length: 99, Replace: "x"}})
	if err == nil {
		t.Fatal("expected out-of-range error")
	}
}

func TestApply_NoTrailingNewline(t *testing.T) {
	src := "no newline"
	edits := []katas.FixEdit{{Line: 1, Column: 1, Length: 2, Replace: "NO"}}
	got, err := Apply(src, edits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "NO newline"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestDiff_Empty(t *testing.T) {
	got, err := Diff("f.zsh", "hello\n", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty diff for no edits, got %q", got)
	}
}

func TestDiff_SimpleReplacement(t *testing.T) {
	src := "echo hi\n"
	edits := []katas.FixEdit{{Line: 1, Column: 1, Length: 4, Replace: "print"}}
	got, err := Diff("f.zsh", src, edits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantHeaders := []string{"--- f.zsh", "+++ f.zsh (fixed)"}
	for _, h := range wantHeaders {
		if !strings.Contains(got, h) {
			t.Errorf("diff missing header %q:\n%s", h, got)
		}
	}
	if !strings.Contains(got, "-echo hi") {
		t.Errorf("diff missing old line:\n%s", got)
	}
	if !strings.Contains(got, "+print hi") {
		t.Errorf("diff missing new line:\n%s", got)
	}
}

func TestDiff_NoChange(t *testing.T) {
	// Edit that replaces with the same text — net zero.
	src := "hi\n"
	edits := []katas.FixEdit{{Line: 1, Column: 1, Length: 2, Replace: "hi"}}
	got, err := Diff("f.zsh", src, edits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty diff for no-op edit, got %q", got)
	}
}
