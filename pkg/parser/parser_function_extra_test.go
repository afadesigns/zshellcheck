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

// Zsh function names can start with `-` (e.g.
// `function -coreutils-alias-setup { … }`). zsh-utils uses this
// pattern for internal-only helpers.
func TestParseFunctionDashPrefixedName(t *testing.T) {
	parseSourceClean(t, "function -coreutils-alias-setup { :; }\n")
}

// Brace-form `if X { } elif Y { }` chain. zinit relies on this Zsh
// shorthand. Previously the parser handled `if X { } else { }` only.
func TestParseBraceFormIfElif(t *testing.T) {
	parseSourceClean(t, "if [[ 1 ]] {\n  echo a\n} elif [[ 2 ]] {\n  echo b\n}\n")
}

func TestParseBraceFormIfElifElse(t *testing.T) {
	parseSourceClean(t, "if [[ 1 ]] {\n  echo a\n} elif [[ 2 ]] {\n  echo b\n} else {\n  echo c\n}\n")
}

// `IDENT+=value` immediately after a brace-form if. Previously failed
// because ParseProgram unconditionally advanced past the head token
// after parseStatement returned, despite the consumedBraceTerminator
// flag indicating the brace-form had already advanced.
func TestParseBraceFormIfFollowedByPlusEq(t *testing.T) {
	parseSourceClean(t, "if [[ 1 ]] {\n  echo a\n}\nx+=1\n")
}

// `((1)); (( y ))` — semicolon between top-level statements must not
// over-advance. Previously ParseProgram consumed the `;` via parseStatement,
// then advanced again, swallowing the second `((`.
func TestParseTopLevelSemicolonChain(t *testing.T) {
	parseSourceClean(t, "((1)); (( y += 1 ))\n")
}

func TestParseTopLevelSemicolonAfterIf(t *testing.T) {
	parseSourceClean(t, "if [[ 1 ]] { :; }; (( x = 1 ))\n")
}

// Nested case inside a case clause body. Previously the inner `esac`
// was treated as the outer body's terminator, collapsing the outer
// case prematurely. zsh-vi-mode uses this pattern in its
// range-handler dispatch.
func TestParseNestedCase(t *testing.T) {
	src := "case $x in\n" +
		"  a)\n" +
		"    case $y in\n" +
		"      b) :;;\n" +
		"    esac\n" +
		"    :\n" +
		"    ;;\n" +
		"  c) :;;\n" +
		"esac\n"
	parseSourceClean(t, src)
}

// Pipeline tail after case must still work after the nested-case fix.
func TestParseCasePipelineTail(t *testing.T) {
	parseSourceClean(t, "case $x in a) echo a;; esac | tr a-z A-Z\n")
}
