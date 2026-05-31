// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
)

func TestParseForLoopArithmeticConditionOnly(t *testing.T) {
	parseSourceClean(t, "for ((i=0; i<3; )) do echo $i; done\n")
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
