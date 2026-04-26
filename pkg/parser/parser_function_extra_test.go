// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import "testing"

func TestParseFunctionLiteralDollarBraceName(t *testing.T) {
	parseSourceClean(t, "function ${=X} { echo hi; }\n")
}

func TestParseFunctionLiteralCompositeName(t *testing.T) {
	parseSourceClean(t, "function name\"${1:-}\"suffix() { echo hi; }\n")
}

func TestParseFunctionLiteralMultiNames(t *testing.T) {
	parseSourceClean(t, "function a b c { echo hi; }\n")
}

func TestParseFunctionLiteralNoParens(t *testing.T) {
	parseSourceClean(t, "function name { echo hi; }\n")
}

func TestParseFunctionLiteralEmptyBody(t *testing.T) {
	parseSourceClean(t, "function name() { }\n")
}

func TestParseGroupedExpressionMultiword(t *testing.T) {
	parseSourceClean(t, "x=(a b c d e)\n")
}

func TestParseGroupedExpressionEmpty(t *testing.T) {
	parseSourceClean(t, "x=()\n")
}

func TestParseIfStatementElif(t *testing.T) {
	parseSourceClean(t, "if true; then echo a; elif false; then echo b; else echo c; fi\n")
}

func TestParseIfStatementNested(t *testing.T) {
	parseSourceClean(t, "if true; then if true; then echo nested; fi; fi\n")
}

func TestParseForLoopBraceStyle(t *testing.T) {
	parseSourceClean(t, "for f in *.zsh; { echo $f }\n")
}

func TestParseForLoopArithmeticHeader(t *testing.T) {
	parseSourceClean(t, "for ((i=0; i<3; i++)); do echo $i; done\n")
}
