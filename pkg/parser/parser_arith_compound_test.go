// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import "testing"

// TestParseArithCompoundDivAssign verifies that (( n /= 1024 )) parses
// without errors — previously the lexer emitted SLASH+ASSIGN instead of
// the single PLUSEQ compound-assign token the parser expects.
func TestParseArithCompoundDivAssign(t *testing.T) {
	parseSourceClean(t, "(( n /= 1024 ))\n")
}

// TestParseArithCompoundBitwiseAndAssign verifies that (( n &= 1 ))
// parses without errors.
func TestParseArithCompoundBitwiseAndAssign(t *testing.T) {
	parseSourceClean(t, "(( n &= 1 ))\n")
}

// TestParseArithCompoundBitwiseOrAssign verifies that (( n |= 1 ))
// parses without errors.
func TestParseArithCompoundBitwiseOrAssign(t *testing.T) {
	parseSourceClean(t, "(( n |= 1 ))\n")
}

// TestParseArithCompoundCaretAssign verifies that (( n ^= 1 )) parses
// without errors.
func TestParseArithCompoundCaretAssign(t *testing.T) {
	parseSourceClean(t, "(( n ^= 1 ))\n")
}

// TestParseArithCompoundLeftShiftAssign verifies that (( n <<= 1 ))
// parses without errors — previously the lexer triggered heredoc
// consumption on the << and the remaining = caused a parse error.
func TestParseArithCompoundLeftShiftAssign(t *testing.T) {
	parseSourceClean(t, "(( n <<= 1 ))\n")
}

// TestParseArithCompoundRightShiftAssign verifies that (( n >>= 1 ))
// parses without errors.
func TestParseArithCompoundRightShiftAssign(t *testing.T) {
	parseSourceClean(t, "(( n >>= 1 ))\n")
}

// TestParseArithCompoundPowerAssign verifies that (( n **= 2 )) parses
// without errors — previously the lexer split ** into two ASTERISK
// tokens and **= into ASTERISK+PLUSEQ, producing a type mismatch.
func TestParseArithCompoundPowerAssign(t *testing.T) {
	parseSourceClean(t, "(( n **= 2 ))\n")
}

// TestParseArithDollarBraceOperand verifies that a parameter expansion
// used as part of an arithmetic operand name parses without errors.
// Previously `PFX_${state}_SFX` caused "expected )), got ${" because
// parseIdentifier did not absorb DollarLbrace tokens in arithmetic.
func TestParseArithDollarBraceOperand(t *testing.T) {
	parseSourceClean(t, "(( x >= PFX_${state}_SFX ))\n")
}

// TestParseArithDollarBraceSimple verifies a bare ${var} in arithmetic.
func TestParseArithDollarBraceSimple(t *testing.T) {
	parseSourceClean(t, "(( x = ${y} + 1 ))\n")
}

// `#` inside `$(( … ))` arithmetic expansion is an operator (the
// char-code prefix `##A` or the positional-arg count `#`), not a comment
// opener, even with a preceding space. The lexer marks the inner `(` of
// `$((` so inArithmetic() holds. Issue #1361.
func TestParseArithDollarParenHashNotComment(t *testing.T) {
	parseSourceClean(t, "x=$(( ##A ))\n")
	parseSourceClean(t, "y=$(( # > 0 ))\n")
}

// A bare `?` in arithmetic prefix position is the `$?` special parameter,
// not the ternary operator. `(( ? == 0 ))` must read as `$? == 0`; the
// operator-followed form used to error "no prefix for ==". The ternary
// `?` (infix) is unaffected. Issue #1378 (used by the Pure prompt).
func TestParseArithBareQuestionBeforeOperator(t *testing.T) {
	parseSourceClean(t, "x=$((? == 0))\n")
	parseSourceClean(t, "y=$(( ? != 0 ))\n")
	parseSourceClean(t, "z=$(( ? + 1 ))\n")
	parseSourceClean(t, "t=$(( cond ? 2 : 3 ))\n")
	parseSourceClean(t, "n=$(( a ? b : c ? d : e ))\n")
}

// TestParseArithAllCompoundOpsRegression is a stress test exercising all
// compound-assign operators in a single arithmetic block.
func TestParseArithAllCompoundOpsRegression(t *testing.T) {
	cases := []string{
		"(( n += 1 ))\n",
		"(( n -= 1 ))\n",
		"(( n *= 2 ))\n",
		"(( n /= 2 ))\n",
		"(( n %= 3 ))\n",
		"(( n &= 1 ))\n",
		"(( n |= 1 ))\n",
		"(( n ^= 1 ))\n",
		"(( n <<= 1 ))\n",
		"(( n >>= 1 ))\n",
		"(( n **= 2 ))\n",
		"(( n = 1 ))\n",
	}
	for _, src := range cases {
		t.Run(src, func(t *testing.T) {
			parseSourceClean(t, src)
		})
	}
}
