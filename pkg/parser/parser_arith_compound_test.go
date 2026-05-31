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
