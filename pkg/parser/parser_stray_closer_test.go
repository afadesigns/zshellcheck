// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import (
	"strings"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
)

// parseStrayExpectErr parses src and fails unless exactly one error is
// reported and it names the unmatched closer.
func parseStrayExpectErr(t *testing.T, src, closer string) {
	t.Helper()
	p := New(lexer.New(src))
	_ = p.ParseProgram()
	errs := p.Errors()
	if len(errs) != 1 {
		t.Fatalf("want 1 error for %q, got %d: %v", src, len(errs), errs)
	}
	if !strings.Contains(errs[0], "unexpected `"+closer+"`") {
		t.Fatalf("error for %q does not name closer %q: %s", src, closer, errs[0])
	}
}

// An orphan compound closer at the top level is unmatched — Zsh rejects
// it. The parser records one error and recovers (issue #1362).
func TestParseStrayCloserTopLevel(t *testing.T) {
	parseStrayExpectErr(t, "done\n", "done")
	parseStrayExpectErr(t, "fi\n", "fi")
	parseStrayExpectErr(t, "esac\n", "esac")
	parseStrayExpectErr(t, "for i in a; do echo; done\ndone\n", "done")
	parseStrayExpectErr(t, "if x; then y; fi\nfi\n", "fi")
	parseStrayExpectErr(t, "echo ok\nesac\n", "esac")
}

// A closer that legitimately closes its own compound, however nested,
// must not be flagged.
func TestParseNestedClosersClean(t *testing.T) {
	parseClean(t, "while x; do for i in a; do echo; done; done\n")
	parseClean(t, "if x; then for i in a; do echo; done; fi\n")
	parseClean(t, "case $x in a) echo a ;; esac\n")
}

// A closer keyword used as a plain command word is not a stray closer.
func TestParseCloserAsWordClean(t *testing.T) {
	parseClean(t, "echo done\n")
	parseClean(t, "a && echo done\n")
	parseClean(t, "echo done\necho fi\n")
}

// A command followed only by redirects parses as one statement: the
// redirect chain must not orphan a target word into a bogus second
// statement. Covers bare (`>>`, `>&`, `<&`), the combined `&>`, and the
// `>/dev/null 2>&1` FD-prefixed tail.
func TestParseRedirectChains(t *testing.T) {
	for _, src := range []string{
		"cmd >out\n",
		"cmd >> log\n",
		"cmd >& out\n",
		"cmd <& 3\n",
		"cmd <in\n",
		"cmd >/dev/null 2>&1\n",
		"cmd >out 2>&1\n",
		"cmd >out >>log\n",
		"cmd 2>&1 >out\n",
		"cmd arg >out 2>&1\n",
	} {
		p := New(lexer.New(src))
		prog := p.ParseProgram()
		if errs := p.Errors(); len(errs) != 0 {
			t.Fatalf("unexpected errors for %q: %v", src, errs)
		}
		if got := len(prog.Statements); got != 1 {
			t.Fatalf("want 1 statement for %q, got %d", src, got)
		}
	}
}

// `&>file` is the combined stdout+stderr redirect, lexed as GTAMP. It
// must not background the command (`&`) and orphan the `>`.
func TestParseAmpGtRedirect(t *testing.T) {
	for _, src := range []string{
		"cmd &>/dev/null\n",
		"cmd &> out\n",
		"zpty -t $worker &>/dev/null || return 1\n",
		"if { cmd } &>/dev/null; then x; fi\n",
	} {
		p := New(lexer.New(src))
		prog := p.ParseProgram()
		if errs := p.Errors(); len(errs) != 0 {
			t.Fatalf("unexpected errors for %q: %v", src, errs)
		}
		if got := len(prog.Statements); got != 1 {
			t.Fatalf("want 1 statement for %q, got %d", src, got)
		}
	}
}

// An expression-led pipeline (`$cmd && echo arg`) gathers the right-hand
// command's arguments instead of stranding them as a separate statement.
func TestParseExpressionLedPipelineArgs(t *testing.T) {
	for _, src := range []string{
		"$cmd && echo hello\n",
		"$cmd && echo done\n",
		"${x} && echo hi there\n",
		"`foo` && bar baz qux\n",
		"$cmd && echo a || echo b\n",
	} {
		p := New(lexer.New(src))
		prog := p.ParseProgram()
		if errs := p.Errors(); len(errs) != 0 {
			t.Fatalf("unexpected errors for %q: %v", src, errs)
		}
		if got := len(prog.Statements); got != 1 {
			t.Fatalf("want 1 statement for %q, got %d", src, got)
		}
	}
}

// A leading `{` in an if/while condition is a brace-group condition, not
// the brace-form body opener (`if { cmd }; then …; fi`).
func TestParseBraceGroupCondition(t *testing.T) {
	parseClean(t, "if { cmd }; then x; fi\n")
	parseClean(t, "if { a || b }; then x; fi\n")
	parseClean(t, "while { c }; do x; done\n")
}

// Zsh's `{ try } always { always-list }` parses as one block; the
// always-list is not stranded as a bogus `always { … }` command.
func TestParseAlwaysBlock(t *testing.T) {
	for _, src := range []string{
		"{ echo try } always { echo cleanup }\n",
		"{ false && true } always { echo done }\n",
		"{ a } always { b } always { c }\n",
	} {
		p := New(lexer.New(src))
		prog := p.ParseProgram()
		if errs := p.Errors(); len(errs) != 0 {
			t.Fatalf("unexpected errors for %q: %v", src, errs)
		}
		if got := len(prog.Statements); got != 1 {
			t.Fatalf("want 1 statement for %q, got %d", src, got)
		}
		block, ok := prog.Statements[0].(*ast.BlockStatement)
		if !ok {
			t.Fatalf("want *ast.BlockStatement for %q, got %T", src, prog.Statements[0])
		}
		if len(block.Statements) < 2 {
			t.Fatalf("always-list not folded into block for %q: %d statements", src, len(block.Statements))
		}
	}
}

// A bare `always` that is not followed by a brace block is consumed but
// not folded; the parser must not panic or error on the degenerate form.
func TestParseAlwaysWithoutBlock(t *testing.T) {
	for _, src := range []string{
		"{ echo a } always\n",
		"{ echo a } always cmd\n",
	} {
		p := New(lexer.New(src))
		_ = p.ParseProgram()
		if errs := p.Errors(); len(errs) != 0 {
			t.Fatalf("unexpected errors for %q: %v", src, errs)
		}
	}
}
