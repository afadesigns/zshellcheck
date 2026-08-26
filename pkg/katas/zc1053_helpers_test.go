// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

// ident builds a bare-word argument.
func zc1053Ident(v string) ast.Expression { return &ast.Identifier{Value: v} }

// str builds a quoted or unquoted string-literal argument. StringLiteral
// keeps its surrounding quotes, which is how the kata tells a redirection
// from a search pattern.
func zc1053Str(v string) ast.Expression { return &ast.StringLiteral{Value: v} }

// TestZC1053QuietFlagScan covers the option-cluster walk: a real quiet flag
// enables silence, an option that consumes a value swallows a following `q`,
// `--` ends option scanning, and a quoted word is an operand.
func TestZC1053QuietFlagScan(t *testing.T) {
	cases := []struct {
		name string
		args []ast.Expression
		want bool
	}{
		{"short quiet", []ast.Expression{zc1053Ident("-q")}, true},
		{"cluster quiet", []ast.Expression{zc1053Ident("-sq")}, true},
		{"long quiet", []ast.Expression{zc1053Ident("--quiet")}, true},
		{"long silent", []ast.Expression{zc1053Ident("--silent")}, true},
		{"other long flag", []ast.Expression{zc1053Ident("--color=auto")}, false},
		{"plain cluster", []ast.Expression{zc1053Ident("-i")}, false},
		{"value flag swallows q", []ast.Expression{zc1053Ident("-fq")}, false},
		{"separate value is skipped", []ast.Expression{zc1053Ident("-e"), zc1053Ident("q")}, false},
		{"quoted operand", []ast.Expression{zc1053Str(`"-quiet"`)}, false},
		{"after double dash", []ast.Expression{zc1053Ident("--"), zc1053Ident("-q")}, false},
		{"bare dash", []ast.Expression{zc1053Ident("-")}, false},
		{"no flags", []ast.Expression{zc1053Ident("pattern")}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := zc1053HasQuietFlag(tc.args); got != tc.want {
				t.Errorf("zc1053HasQuietFlag(%s) = %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}

// TestZC1053DevNullRedirect covers every shape a stdout redirection reaches
// the kata in, plus the forms that only look like one and must still report.
func TestZC1053DevNullRedirect(t *testing.T) {
	concat := func(parts ...ast.Expression) ast.Expression {
		return &ast.ConcatenatedExpression{Parts: parts}
	}
	cases := []struct {
		name string
		args []ast.Expression
		want bool
	}{
		{"glued", []ast.Expression{concat(zc1053Str(">"), zc1053Ident("/dev/null"))}, true},
		{"spaced", []ast.Expression{zc1053Str(">"), zc1053Ident("/dev/null")}, true},
		{"append", []ast.Expression{concat(zc1053Str(">>"), zc1053Ident("/dev/null"))}, true},
		{"both streams", []ast.Expression{concat(zc1053Str("&>"), zc1053Ident("/dev/null"))}, true},
		{"explicit stdout fd", []ast.Expression{concat(zc1053Str("1>"), zc1053Ident("/dev/null"))}, true},
		{"quoted target", []ast.Expression{concat(zc1053Str(">"), zc1053Str(`"/dev/null"`))}, true},
		{"stderr only", []ast.Expression{concat(zc1053Str("2>"), zc1053Ident("/dev/null"))}, false},
		{"look-alike path", []ast.Expression{concat(zc1053Str(">"), zc1053Ident("/dev/null.txt"))}, false},
		{"quoted redirect is a pattern", []ast.Expression{zc1053Str(`'>/dev/null'`)}, false},
		{"quoted concat head", []ast.Expression{concat(zc1053Str(`'>'`), zc1053Ident("/dev/null"))}, false},
		{"dangling operator", []ast.Expression{zc1053Str(">")}, false},
		{"plain word", []ast.Expression{zc1053Ident("pattern")}, false},
		{"fd duplication", []ast.Expression{concat(zc1053Str("2>&"), &ast.IntegerLiteral{Value: 1})}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := zc1053HasDevNullRedirect(tc.args); got != tc.want {
				t.Errorf("zc1053HasDevNullRedirect(%s) = %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}

// TestZC1053SplitRedirect pins the operator/target split, including the
// file-descriptor prefix and the non-redirection case.
func TestZC1053SplitRedirect(t *testing.T) {
	cases := map[string][2]string{
		">/dev/null":  {">", "/dev/null"},
		"1>/dev/null": {"1>", "/dev/null"},
		"2>/dev/null": {"2>", "/dev/null"},
		"&>/dev/null": {"&>", "/dev/null"},
		">|":          {">|", ""},
		"pattern":     {"", ""},
		"":            {"", ""},
	}
	for word, want := range cases {
		op, target := zc1053SplitRedirect(word)
		if op != want[0] || target != want[1] {
			t.Errorf("zc1053SplitRedirect(%q) = (%q, %q), want (%q, %q)", word, op, target, want[0], want[1])
		}
	}
}

// TestZC1053RedirectsStdout pins which operators carry stdout.
func TestZC1053RedirectsStdout(t *testing.T) {
	for _, op := range []string{">", ">>", "&>", ">&", ">|", "1>", "1>>"} {
		if !zc1053RedirectsStdout(op) {
			t.Errorf("zc1053RedirectsStdout(%q) = false, want true", op)
		}
	}
	for _, op := range []string{"2>", "2>&", "3>", "<", ""} {
		if zc1053RedirectsStdout(op) {
			t.Errorf("zc1053RedirectsStdout(%q) = true, want false", op)
		}
	}
}

// TestZC1053WordHelpers covers quote detection and stripping.
func TestZC1053WordHelpers(t *testing.T) {
	if word, quoted := zc1053ArgWord(zc1053Str(`'x'`)); !quoted || word != `'x'` {
		t.Errorf("zc1053ArgWord(quoted) = (%q, %v), want ('x', true)", word, quoted)
	}
	if word, quoted := zc1053ArgWord(zc1053Ident("plain")); quoted || word != "plain" {
		t.Errorf("zc1053ArgWord(identifier) = (%q, %v), want (plain, false)", word, quoted)
	}
	if got := zc1053UnquoteWord(`"/dev/null"`); got != "/dev/null" {
		t.Errorf("zc1053UnquoteWord(double) = %q, want /dev/null", got)
	}
	if got := zc1053UnquoteWord(`'/dev/null'`); got != "/dev/null" {
		t.Errorf("zc1053UnquoteWord(single) = %q, want /dev/null", got)
	}
	if got := zc1053UnquoteWord("/dev/null"); got != "/dev/null" {
		t.Errorf("zc1053UnquoteWord(bare) = %q, want /dev/null", got)
	}
	if zc1053IsQuoted("") {
		t.Error("zc1053IsQuoted(empty) = true, want false")
	}
}

// TestZC1053NonIdentifierCommand covers the guard for a command whose name is
// not a bare word, such as one invoked through a variable.
func TestZC1053NonIdentifierCommand(t *testing.T) {
	violations := []Violation{}
	cmd := &ast.SimpleCommand{Name: &ast.StringLiteral{Value: `"grep"`}}
	checkCommandZC1053(cmd, false, &violations)
	if len(violations) != 0 {
		t.Errorf("expected no violation for a non-identifier command name, got %d", len(violations))
	}
	// A silenced context reports nothing regardless of the command.
	checkCommandZC1053(&ast.SimpleCommand{Name: &ast.Identifier{Value: "grep"}}, true, &violations)
	if len(violations) != 0 {
		t.Errorf("expected no violation in a silenced context, got %d", len(violations))
	}
}

// TestZC1053SilencesStdout covers the redirection-node path, which a
// subshell or brace group carrying the redirect still reaches.
func TestZC1053SilencesStdout(t *testing.T) {
	devNull := &ast.Identifier{Value: "/dev/null"}
	for _, op := range []string{">", ">>", "&>"} {
		red := &ast.Redirection{Operator: op, Right: devNull}
		if !zc1053SilencesStdout(red) {
			t.Errorf("zc1053SilencesStdout(%q /dev/null) = false, want true", op)
		}
	}
	if zc1053SilencesStdout(&ast.Redirection{Operator: ">", Right: &ast.Identifier{Value: "/tmp/log"}}) {
		t.Error("zc1053SilencesStdout(> /tmp/log) = true, want false")
	}
	if zc1053SilencesStdout(&ast.Redirection{Operator: "<", Right: devNull}) {
		t.Error("zc1053SilencesStdout(< /dev/null) = true, want false")
	}
}

// TestZC1293LastExprArg pins where the rewritten conditional ends. A
// trailing redirection belongs outside the brackets, and a quoted word that
// merely looks like one is part of the expression.
func TestZC1293LastExprArg(t *testing.T) {
	concat := func(parts ...ast.Expression) ast.Expression {
		return &ast.ConcatenatedExpression{Parts: parts}
	}
	cases := []struct {
		name string
		args []ast.Expression
		want int
	}{
		{"no redirection", []ast.Expression{zc1053Ident("-f"), zc1053Ident("file")}, 1},
		{
			"attached redirection",
			[]ast.Expression{zc1053Ident("-f"), zc1053Ident("f"), concat(zc1053Str("2>"), zc1053Ident("/dev/null"))},
			1,
		},
		{
			"spaced redirection",
			[]ast.Expression{zc1053Ident("-f"), zc1053Ident("f"), zc1053Str(">"), zc1053Ident("/dev/null")},
			1,
		},
		{
			// A quoted word is an operand even when it looks like an operator.
			"quoted operator is an operand",
			[]ast.Expression{zc1053Ident("-n"), zc1053Str(`">"`)},
			1,
		},
		{"only a redirection", []ast.Expression{concat(zc1053Str(">"), zc1053Ident("/dev/null"))}, -1},
		{"no arguments", []ast.Expression{}, -1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := zc1293LastExprArg(tc.args); got != tc.want {
				t.Errorf("zc1293LastExprArg(%s) = %d, want %d", tc.name, got, tc.want)
			}
		})
	}
}
