// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
)

// TestParseEnvPrefixAssignmentFlag locks the #1332 fix: an assignment
// that prefixes a command on the same line is marked EnvPrefix, while a
// standalone assignment (or one ended by `;`) is not.
func TestParseEnvPrefixAssignmentFlag(t *testing.T) {
	cases := []struct {
		src  string
		want bool
	}{
		{"DEBUG=true echo foo\n", true},   // inline env-var prefix
		{"x+=1 mycmd\n", true},            // `+=` prefix form
		{"DEBUG=true\n", false},           // standalone assignment
		{"DEBUG=true; echo foo\n", false}, // `;` ends the assignment
		{"echo foo\n", false},             // not an assignment at all
		{"a[0]=1 echo foo\n", false},      // indexed LHS is not a plain name
	}
	for _, tc := range cases {
		prog := New(lexer.New(tc.src)).ParseProgram()
		es, ok := prog.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("%q: statement 0 is %T, want *ast.ExpressionStatement", tc.src, prog.Statements[0])
		}
		if es.EnvPrefix != tc.want {
			t.Errorf("%q: EnvPrefix = %v, want %v", tc.src, es.EnvPrefix, tc.want)
		}
	}
}

func TestParseForLoopArithmeticConditionOnly(t *testing.T) {
	parseSourceClean(t, "for ((i=0; i<3; )) do echo $i; done\n")
}

// A POSIX bracket class in a for-in word (`a[[:alpha:]]`) is a glob, not
// an array subscript. parseIndexExpression used to drain forward looking
// for a `]` the mis-parse had already consumed, swallowing `; do … done`
// and the rest of the input. Issue #1376 (blocks Powerlevel10k).
func TestParseForInBracketClass(t *testing.T) {
	parseSourceClean(t, "for x in a[[:alpha:]]; do echo \"$x\"; done\n")
	parseSourceClean(t, "for f in /etc/*[[:digit:]]; do :; done\n")
}

func TestParseForLoopArithmeticCommaInit(t *testing.T) {
	parseSourceClean(t, "for ((i=0, j=10; i<j; i++)) do echo $i; done\n")
}

func TestParseForLoopArithmeticCommaPost(t *testing.T) {
	parseSourceClean(t, "for ((i=0; i<10; i++, j--)) do echo $i; done\n")
}

func TestParseForLoopArithmeticCommaBoth(t *testing.T) {
	parseSourceClean(t, "for ((i=0, j=10; i<j; i++, j--)) do echo $i $j; done\n")
}

// The comma operator must chain into the slot expression, not be
// dropped: stmt.Init has to carry both `i=0` and `j=10`.
func TestParseForLoopArithmeticCommaChains(t *testing.T) {
	prog := New(lexer.New("for ((i=0, j=10; i<j; )) do echo $i; done\n")).ParseProgram()
	stmt, ok := prog.Statements[0].(*ast.ForLoopStatement)
	if !ok {
		t.Fatalf("Statements[0] is not *ast.ForLoopStatement; got %T", prog.Statements[0])
	}
	infix, ok := stmt.Init.(*ast.InfixExpression)
	if !ok || infix.Operator != "," {
		t.Fatalf("Init is not a comma-chained InfixExpression; got %T", stmt.Init)
	}
}

func TestParseForLoopArithmeticInitOnly(t *testing.T) {
	parseSourceClean(t, "for ((i=0; ; )) do break; done\n")
}

func TestParseForLoopMultiVariable(t *testing.T) {
	parseSourceClean(t, "for k v in a 1 b 2; do echo $k $v; done\n")
}

func TestParseForLoopShortForm(t *testing.T) {
	parseSourceClean(t, "for f (a b c) echo $f\n")
}

func TestParseForLoopImplicitList(t *testing.T) {
	parseSourceClean(t, "for k v w; do echo $k; done\n")
}

func TestParseForLoopNumericName(t *testing.T) {
	parseSourceClean(t, "for 1 in a b c; do echo $1; done\n")
}

func TestParseIfStatementInline(t *testing.T) {
	parseSourceClean(t, "if [ -f f ]; then echo y; fi\n")
}

func TestParseIfStatementWithCommandSubstitution(t *testing.T) {
	parseSourceClean(t, "if [[ $(date +%H) -lt 12 ]]; then echo morning; fi\n")
}

func TestParseSubshellStatement(t *testing.T) {
	parseSourceClean(t, "(cd /tmp && rm -rf foo)\n")
}

func TestParseSubshellWithPipeline(t *testing.T) {
	parseSourceClean(t, "(echo a; echo b) | sort\n")
}

func TestParseDeclarationValueArray(t *testing.T) {
	parseSourceClean(t, "x=(1 2 3 4)\n")
}

func TestParseDeclarationValueAssoc(t *testing.T) {
	parseSourceClean(t, "typeset -A m=(k1 v1 k2 v2)\n")
}
