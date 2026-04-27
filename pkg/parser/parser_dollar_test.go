// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/lexer"
)

func TestParseBareDollarBeforeSemicolon(t *testing.T) {
	parseSourceClean(t, "echo $;\n")
}

func TestParseBareDollarBeforePipe(t *testing.T) {
	parseSourceClean(t, "echo $ | wc\n")
}

func TestParseBareDollarBeforeRparen(t *testing.T) {
	parseSourceClean(t, "x=( $ )\n")
}

func TestParseDollarBracketArithmetic(t *testing.T) {
	parseSourceClean(t, "echo $[1+2]\n")
}

func TestParseDollarHashLength(t *testing.T) {
	parseSourceClean(t, "echo $#name\n")
}

// $#1 — length of positional $1. Used inside `for ((i=1;i<=$#1;i++))`
// loops in zsh-vi-mode. Previously the parser rejected the trailing
// digit because parseDollarSpecialOp only accepted IDENT after `#`.
func TestParseDollarHashPositional(t *testing.T) {
	parseSourceClean(t, "for ((i=1;i<=$#1;i++)); do :; done\n")
}

func TestParseDollarHashStar(t *testing.T) {
	parseSourceClean(t, "echo $#*\n")
}

func TestParseDollarHashQuestion(t *testing.T) {
	parseSourceClean(t, "echo $#?\n")
}

func TestParseDollarInteger(t *testing.T) {
	parseSourceClean(t, "echo $1\n")
}

func TestParseDollarBang(t *testing.T) {
	parseSourceClean(t, "echo $!\n")
}

func TestParseDollarMinus(t *testing.T) {
	parseSourceClean(t, "echo $-\n")
}

func TestParseDollarStar(t *testing.T) {
	parseSourceClean(t, "echo $*\n")
}

func TestParseDollarPlusName(t *testing.T) {
	parseSourceClean(t, "(( $+name ))\n")
}

func TestParseDollarPlusNameSubscript(t *testing.T) {
	parseSourceClean(t, "(( $+name[key] ))\n")
}

// $((`[##N]` …)) — Zsh arithmetic radix prefix that prints the result
// in a non-decimal base. zsh-vi-mode uses `h=$(([##16]$h+1))` to
// generate hex-encoded keymap names. The parser previously failed
// because LBRACKET inside `((` is registered as parseSingleCommand
// (command-position `[ ... ]`), not as a radix opener.
func TestParseArithRadixHashHash(t *testing.T) {
	parseSourceClean(t, "h=$(([##16]$h+1))\n")
}

func TestParseArithRadixHashOne(t *testing.T) {
	parseSourceClean(t, "echo $(([#8]7))\n")
}

func TestParseArithRadixHashOnly(t *testing.T) {
	parseSourceClean(t, "echo $(([#]42))\n")
}

func TestParseDoubleParenRadix(t *testing.T) {
	parseSourceClean(t, "(([##16]$x+1))\n")
}

// Arithmetic comparison with hex/binary literal prefix and parameter
// expansion. zsh-vi-mode line 2734: `(( $number == 0x${(l:15::f:)} ))`.
// The lexer extension keeps `0x…` digits as one INT token; the parser
// then treats `0x` followed by `${…}` as INT(0) + IDENT(x) + DollarLbrace
// (degenerate Zsh string concat) and recovers cleanly.
func TestParseArithHexWithExpansion(t *testing.T) {
	parseSourceClean(t, "(( $number == 0x${var} ))\n")
}

func TestParseArithBinaryWithExpansion(t *testing.T) {
	parseSourceClean(t, "(( $number == 0b${var} ))\n")
}

func TestParseArithBinaryLiteral(t *testing.T) {
	parseSourceClean(t, "(( $number == 0b101 ))\n")
}

func TestParseArithCustomBaseLiteral(t *testing.T) {
	parseSourceClean(t, "(( x = 16#ff ))\n")
}

// Drive every absorbArithmeticNumberTail branch.
func TestParseArithNumTailVariableConcat(t *testing.T) {
	parseSourceClean(t, "(( a == 0x$h ))\n")
}

func TestParseArithNumTailIntConcat(t *testing.T) {
	parseSourceClean(t, "(( a == 16#1234 ))\n")
}

func TestParseArithNumTailBraceVar(t *testing.T) {
	parseSourceClean(t, "(( a == 0b${(l:8::0:)bits} ))\n")
}

// Drive every consumeArithmeticRadixPrefix branch.
func TestParseArithNoRadixIdentStart(t *testing.T) {
	parseSourceClean(t, "(( x + 1 ))\n")
}

func TestParseArithRadixUnclosed(t *testing.T) {
	// Recovery path: malformed radix prefix without `]` exits at EOF.
	src := "(( [#16\n"
	p := New(lexer.New(src))
	_ = p.ParseProgram()
	if len(p.Errors()) == 0 {
		t.Errorf("expected parser errors for unclosed radix; got none")
	}
}

// Drive every peekIsHashLengthOperand branch.
func TestParseDollarHashHash(t *testing.T) {
	parseSourceClean(t, "echo $##\n")
}

func TestParseDollarHashWithSpace(t *testing.T) {
	// Space after $# means it's not a length op — falls through.
	parseSourceClean(t, "echo $# arg\n")
}

// `(( # ))` — `#` standalone is the count of positional args inside
// arithmetic. zimfw uses `(( ! # ))` for "no args" and
// `(( # > 0 ))` for "have args".
func TestParseArithBareHashCount(t *testing.T) {
	parseSourceClean(t, "(( ! # ))\n")
}

func TestParseArithBareHashCompare(t *testing.T) {
	parseSourceClean(t, "(( # > 0 ))\n")
}

// Backtick comment-suppression inside arithmetic must not eat the
// closing `))`. `#` outside arithmetic is still a comment opener.
func TestParseHashIsCommentOutsideArith(t *testing.T) {
	parseSourceClean(t, "echo before # this is a comment\n")
}

// `$+3` — Zsh existence test for positional `$3` inside arithmetic.
// fzf-tab's ls-colors.zsh uses `if (($+3)); then …`.
func TestParseDollarPlusInt(t *testing.T) {
	parseSourceClean(t, "(( $+3 ))\n")
}

// Bitwise `&` and `^` are infix operators inside `((…))`.
// zinit-autoload's mode-flag tests use `(( unpacked[3] & 0x1 ))`.
func TestParseArithBitwiseAnd(t *testing.T) {
	parseSourceClean(t, "(( unpacked[3] & 0x1 ))\n")
}

func TestParseArithBitwiseXor(t *testing.T) {
	parseSourceClean(t, "(( a ^ b ))\n")
}

// `&` outside arithmetic still backgrounds the command.
func TestParseAmpersandStillBackgrounds(t *testing.T) {
	parseSourceClean(t, "sleep 5 &\n")
}

// Comma operator inside arithmetic. zinit chains side effects via
// `(( ++idx, count += val ))`.
func TestParseArithCommaOperator(t *testing.T) {
	parseSourceClean(t, "(( ++a, b += 1 ))\n")
}

func TestParseArithCommaInDollarParen(t *testing.T) {
	parseSourceClean(t, "echo $(( ++a, b += 1 ))\n")
}

func TestParseFuncCallCommasUnaffected(t *testing.T) {
	parseSourceClean(t, "let x = add(a, b, 1)\n")
}
