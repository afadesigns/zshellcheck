// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import "testing"

func TestParseCaseStatementBasic(t *testing.T) {
	parseSourceClean(t, "case $x in a) echo a;; b) echo b;; esac\n")
}

func TestParseCaseStatementMultiplePatterns(t *testing.T) {
	parseSourceClean(t, "case $x in a|b|c) echo abc;; *) echo other;; esac\n")
}

func TestParseCaseStatementParenLabel(t *testing.T) {
	parseSourceClean(t, "case $x in (a) echo a;; (b) echo b;; esac\n")
}

func TestParseCaseStatementGlobAlternation(t *testing.T) {
	parseSourceClean(t, "case $x in (darwin|freebsd)*) echo bsd;; esac\n")
}

func TestParseCaseStatementNested(t *testing.T) {
	parseSourceClean(t, "case x in a) case y in 1) echo nested;; esac;; esac\n")
}

func TestParseCaseStatementEmptyClauses(t *testing.T) {
	parseSourceClean(t, "case $x in a) ;; b) ;; esac\n")
}

func TestParseAnonymousFunction(t *testing.T) {
	parseSourceClean(t, "() { echo anon; }\n")
}

func TestParseShebang(t *testing.T) {
	parseSourceClean(t, "#!/usr/bin/env zsh\necho ok\n")
}

func TestParseDoubleBracketTest(t *testing.T) {
	parseSourceClean(t, "[[ -f file && -r file ]]\n")
}

func TestParseDoubleBracketRegex(t *testing.T) {
	parseSourceClean(t, "[[ $x =~ ^[a-z]+$ ]]\n")
}

// Inside `[[ … ]]`, a leading `[` opens a glob bracket-class
// fragment, not the `[` test-builtin or an array subscript. The
// default LBRACKET prefix used to gobble through the closing `]]`.
func TestParseDoubleBracketGlobBracketClass(t *testing.T) {
	parseSourceClean(t, "[[ $x = [abc]* ]]\n")
}

func TestParseDoubleBracketPosixClass(t *testing.T) {
	parseSourceClean(t, "[[ $x = [[:alnum:]]## ]]\n")
}

func TestParseDoubleBracketNegatedPosixClass(t *testing.T) {
	parseSourceClean(t, "[[ $x = [[:blank:]]##[^[:blank:]]* ]]\n")
}

func TestParseDoubleBracketPosixClassWithCharLiteral(t *testing.T) {
	parseSourceClean(t, "[[ $x = ([[:alpha:]_][[:alnum:]_]#) ]]\n")
}

// Case-clause patterns can carry `[…]` glob bracket-class fragments.
// The default RBRACKET prefix used to invite an INDEX infix, so the
// inner `]` of `[*?]` recursed into a phantom array subscript and
// the outer `)` of the case label became orphaned.
func TestParseCaseClauseGlobBracketPattern(t *testing.T) {
	parseSourceClean(t, "case x in [*?]*|*[^\\\\][*?]*) echo y;; esac\n")
}

// Zsh case-clause fall-through markers `;|` (test next) and `;&`
// (execute next) terminate a clause body the same way `;;` does. The
// lexer fuses them into DSEMI so the parser's case loop honours the
// boundary.
func TestParseCaseClauseFallThroughTestNext(t *testing.T) {
	parseSourceClean(t, "case x in a) echo a;| b) echo b;; esac\n")
}

func TestParseCaseClauseFallThroughExecuteNext(t *testing.T) {
	parseSourceClean(t, "case x in a) echo a;& b) echo b;; esac\n")
}

func TestParseProcessSubstitution(t *testing.T) {
	parseSourceClean(t, "diff <(sort a) <(sort b)\n")
}

func TestParseSelectStatement(t *testing.T) {
	parseSourceClean(t, "select opt in a b c; do echo $opt; break; done\n")
}
