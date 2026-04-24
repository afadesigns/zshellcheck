package fix

import (
	"strings"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
	"github.com/afadesigns/zshellcheck/pkg/parser"
)

// runFix walks source through lex + parse, collects fix edits from the
// shipped registry, applies them, and returns the rewritten source.
// It is used by each kata-level integration test below; keeping it
// centralised means the tests stay small and identical in shape.
func runFix(t *testing.T, source string) string {
	t.Helper()
	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()
	if errs := p.Errors(); len(errs) != 0 {
		t.Fatalf("parse errors: %v", errs)
	}
	var edits []katas.FixEdit
	ast.Walk(program, func(n ast.Node) bool {
		_, es := katas.Registry.CheckAndFix(n, nil, []byte(source))
		edits = append(edits, es...)
		return true
	})
	out, err := Apply(source, edits)
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	return out
}

func TestFixIntegration_ZC1002_Backticks(t *testing.T) {
	src := "result=`ls -la`\n"
	want := "result=$(ls -la)\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1005_Which(t *testing.T) {
	src := "which git\n"
	want := "whence git\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_ZC1092_Echo(t *testing.T) {
	src := "echo hello world\n"
	got := runFix(t, src)
	// ZC1092 replaces just the command name; arguments stay intact.
	if !strings.HasPrefix(got, "print -r --") {
		t.Errorf("expected fix to replace echo, got %q", got)
	}
	if !strings.Contains(got, "hello world") {
		t.Errorf("fix mangled arguments: %q", got)
	}
}

func TestFixIntegration_EchoFlag_LeftAlone(t *testing.T) {
	// ZC1092 Fix skips flagged forms — `echo -n` translates to a
	// different print invocation than `print -r --`.
	src := "echo -n keep\n"
	if got := runFix(t, src); got != src {
		t.Errorf("flagged echo should not be auto-fixed, got %q", got)
	}
}

func TestFixIntegration_NestedKatas_OuterWins(t *testing.T) {
	// Outer backtick-span fix (ZC1002) wraps an inner `which` that
	// ZC1005 would rewrite. The outer edit wins on the first pass —
	// the output preserves `which` inside the new parens. A second
	// `-fix` pass would then convert `which` -> `whence`.
	src := "result=`which git`\n"
	want := "result=$(which git)\n"
	if got := runFix(t, src); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFixIntegration_SecondPass_ResolvesInner(t *testing.T) {
	src := "result=`which git`\n"
	first := runFix(t, src)
	final := runFix(t, first)
	want := "result=$(whence git)\n"
	if final != want {
		t.Errorf("got %q, want %q", final, want)
	}
}
