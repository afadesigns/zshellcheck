// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import "testing"

func TestParseFunctionLiteralKeywordForm(t *testing.T) {
	parseSourceClean(t, "function name() { echo hi; }\n")
}

func TestParseFunctionLiteralBareForm(t *testing.T) {
	parseSourceClean(t, "name() { echo hi; }\n")
}

func TestParseFunctionLiteralKeywordWithoutParens(t *testing.T) {
	parseSourceClean(t, "function greet { echo hi; }\n")
}

func TestParseFunctionLiteralBody(t *testing.T) {
	parseSourceClean(t, "f() { local x=1; echo $x; }\n")
}

// `function` as the right-hand side of an assignment is a literal
// identifier, not a function-definition opener. Without the
// peekIsFunctionDefinitionContinuation guard, parseFunctionLiteral
// errored on the missing `{` body and broke the surrounding elif chain.
func TestParseFunctionKeywordAsAssignmentRhs(t *testing.T) {
	parseSourceClean(t, "REPLY=function\n")
}

func TestParseFunctionKeywordAsAssignmentRhsInElifChain(t *testing.T) {
	parseSourceClean(t, "if (( a )); then REPLY=alias; elif (( b )); then REPLY=function; fi\n")
}
