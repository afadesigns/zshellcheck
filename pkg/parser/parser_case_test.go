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

// Two array assignments inside a subshell. The first `arr=( "x" )`
// closes its array literal on `)`; without setting
// consumedParenTerminator, parseBlockStatement misread that `)` as
// the subshell's terminator and the second assignment never parsed.
func TestParseArrayAssignmentsInsideSubshell(t *testing.T) {
	parseSourceClean(t, "( arr=( \"x\" ); list=( \"y\" ) )\n")
}

func TestParseArrayAssignmentsInsideSubshellNewlineSeparated(t *testing.T) {
	parseSourceClean(t, "(\narr=( \"x\" )\nlist=( \"y\" )\n)\n")
}

// `${#}` is the special parameter "count of positional args", not a
// length operator over a missing subject. Without RBRACE in the
// subjectIsEmpty set, parseArrayAccess advanced past `#` looking for
// a subject and erroring on the closing `}`.
func TestParseDollarBraceHashSpecialParameter(t *testing.T) {
	parseSourceClean(t, "[[ ${#} = 1 ]]\n")
}

// Zsh keywords double as variable names in `${…}` subject position.
// `${(flags)in}` and `${for}` previously errored on the keyword token
// because parseArrayAccessSubject's fallthrough hit
// noPrefixParseFnError on IN / FOR.
func TestParseDollarBraceKeywordSubject(t *testing.T) {
	parseSourceClean(t, "echo ${in} ${for} ${while}\n")
}

func TestParseDollarBraceFlagsKeywordSubject(t *testing.T) {
	parseSourceClean(t, "echo ${(j: :)in}\n")
}

// Zsh shortcut: `if (( cond )) cmd` and `if [[ cond ]] cmd` omit the
// `then`/`fi` pair. Inside `=( … )` proc-sub or `( … )` subshell,
// parseBlockStatement(THEN, LBRACE) absorbed the trailing cmd into
// the cond block and walked past the enclosing terminator.
func TestParseIfShortcutInProcSub(t *testing.T) {
	parseSourceClean(t, "foo =(\n  if (( x )) print y\n)\n")
}

func TestParseIfShortcutDoubleBracketInProcSub(t *testing.T) {
	parseSourceClean(t, "foo =(\n  if [[ z == w ]] echo hi\n)\n")
}

func TestParseIfShortcutInsideFunctionBody(t *testing.T) {
	parseSourceClean(t, "foo() {\n  if (( x )) print y\n}\n")
}

func TestParseProcessSubstitution(t *testing.T) {
	parseSourceClean(t, "diff <(sort a) <(sort b)\n")
}

func TestParseSelectStatement(t *testing.T) {
	parseSourceClean(t, "select opt in a b c; do echo $opt; break; done\n")
}
