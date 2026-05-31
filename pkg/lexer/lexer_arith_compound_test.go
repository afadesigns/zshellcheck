// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package lexer

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/token"
)

// collectArithTokenTypes lexes src and returns the token types between
// the DoubleLparen opener and the DoubleRparen closer, exclusive.
// Used by the arithmetic compound-assign regression tests to check the
// operator token without depending on surrounding frame tokens.
func collectArithTokenTypes(t *testing.T, src string) []token.Type {
	t.Helper()
	toks := tokensFor(t, src)
	var inner []token.Type
	inArith := false
	for _, tok := range toks {
		switch tok.Type {
		case token.DoubleLparen:
			inArith = true
		case token.DoubleRparen, token.EOF:
			return inner
		default:
			if inArith {
				inner = append(inner, tok.Type)
			}
		}
	}
	return inner
}

// TestArithCompoundDivAssign verifies that /= inside (( )) lexes as a
// single PLUSEQ compound-assign token, not SLASH followed by ASSIGN.
func TestArithCompoundDivAssign(t *testing.T) {
	types := collectArithTokenTypes(t, "(( n /= 1024 ))\n")
	// Expect: IDENT('/'), PLUSEQ, INT
	if len(types) != 3 {
		t.Fatalf("/= in arith: got token types %v, want 3 tokens", types)
	}
	if types[1] != token.PLUSEQ {
		t.Errorf("/= in arith: token[1] = %s, want PLUSEQ", types[1])
	}
}

// TestArithCompoundBitwiseAndAssign verifies that &= inside (( )) lexes
// as PLUSEQ, not AMPERSAND followed by ASSIGN.
func TestArithCompoundBitwiseAndAssign(t *testing.T) {
	types := collectArithTokenTypes(t, "(( n &= 1 ))\n")
	if len(types) != 3 {
		t.Fatalf("&= in arith: got token types %v, want 3 tokens", types)
	}
	if types[1] != token.PLUSEQ {
		t.Errorf("&= in arith: token[1] = %s, want PLUSEQ", types[1])
	}
}

// TestArithCompoundBitwiseOrAssign verifies that |= inside (( )) lexes
// as PLUSEQ, not PIPE followed by ASSIGN.
func TestArithCompoundBitwiseOrAssign(t *testing.T) {
	types := collectArithTokenTypes(t, "(( n |= 1 ))\n")
	if len(types) != 3 {
		t.Fatalf("|= in arith: got token types %v, want 3 tokens", types)
	}
	if types[1] != token.PLUSEQ {
		t.Errorf("|= in arith: token[1] = %s, want PLUSEQ", types[1])
	}
}

// TestArithCompoundCaretAssign verifies that ^= inside (( )) lexes as
// PLUSEQ, not CARET followed by ASSIGN.
func TestArithCompoundCaretAssign(t *testing.T) {
	types := collectArithTokenTypes(t, "(( n ^= 1 ))\n")
	if len(types) != 3 {
		t.Fatalf("^= in arith: got token types %v, want 3 tokens", types)
	}
	if types[1] != token.PLUSEQ {
		t.Errorf("^= in arith: token[1] = %s, want PLUSEQ", types[1])
	}
}

// TestArithCompoundLeftShiftAssign verifies that <<= inside (( )) lexes
// as a single PLUSEQ token, not LTLT followed by ASSIGN (heredoc open).
func TestArithCompoundLeftShiftAssign(t *testing.T) {
	types := collectArithTokenTypes(t, "(( n <<= 1 ))\n")
	if len(types) != 3 {
		t.Fatalf("<<= in arith: got token types %v, want 3 tokens", types)
	}
	if types[1] != token.PLUSEQ {
		t.Errorf("<<= in arith: token[1] = %s, want PLUSEQ", types[1])
	}
}

// TestArithCompoundRightShiftAssign verifies that >>= inside (( )) lexes
// as a single PLUSEQ token, not GTGT followed by ASSIGN.
func TestArithCompoundRightShiftAssign(t *testing.T) {
	types := collectArithTokenTypes(t, "(( n >>= 1 ))\n")
	if len(types) != 3 {
		t.Fatalf(">>= in arith: got token types %v, want 3 tokens", types)
	}
	if types[1] != token.PLUSEQ {
		t.Errorf(">>= in arith: token[1] = %s, want PLUSEQ", types[1])
	}
}

// TestArithCompoundPowerAssign verifies that **= inside (( )) lexes as
// a single PLUSEQ token, not ASTERISK followed by PLUSEQ.
func TestArithCompoundPowerAssign(t *testing.T) {
	types := collectArithTokenTypes(t, "(( n **= 2 ))\n")
	if len(types) != 3 {
		t.Fatalf("**= in arith: got token types %v, want 3 tokens", types)
	}
	if types[1] != token.PLUSEQ {
		t.Errorf("**= in arith: token[1] = %s, want PLUSEQ", types[1])
	}
}

// TestNonArithAmpersandBackground verifies that & outside arithmetic
// still lexes as AMPERSAND (background operator), not PLUSEQ.
func TestNonArithAmpersandBackground(t *testing.T) {
	toks := tokensFor(t, "sleep 1 &\n")
	for _, tok := range toks {
		if tok.Type == token.PLUSEQ {
			t.Errorf("& outside arith: got PLUSEQ, want AMPERSAND for background")
		}
	}
}

// TestNonArithPipePipe verifies that | outside arithmetic still lexes
// as PIPE, not PLUSEQ.
func TestNonArithPipePipe(t *testing.T) {
	toks := tokensFor(t, "echo foo | cat\n")
	found := false
	for _, tok := range toks {
		if tok.Type == token.PIPE {
			found = true
		}
		if tok.Type == token.PLUSEQ {
			t.Errorf("| outside arith: got PLUSEQ, want PIPE")
		}
	}
	if !found {
		t.Errorf("| outside arith: PIPE token not found in stream")
	}
}

// TestNonArithGtGtRedirection verifies that >> outside arithmetic still
// lexes as GTGT (append redirection), not PLUSEQ.
func TestNonArithGtGtRedirection(t *testing.T) {
	toks := tokensFor(t, "echo hi >> /dev/null\n")
	found := false
	for _, tok := range toks {
		if tok.Type == token.GTGT {
			found = true
		}
		if tok.Type == token.PLUSEQ {
			t.Errorf(">> outside arith: got PLUSEQ, want GTGT")
		}
	}
	if !found {
		t.Errorf(">> outside arith: GTGT token not found in stream")
	}
}

// TestNonArithLtLtHeredoc verifies that << outside arithmetic still
// lexes as LTLT (heredoc opener), not PLUSEQ.
func TestNonArithLtLtHeredoc(t *testing.T) {
	toks := tokensFor(t, "cat <<EOF\nhello\nEOF\n")
	found := false
	for _, tok := range toks {
		if tok.Type == token.LTLT {
			found = true
		}
		if tok.Type == token.PLUSEQ {
			t.Errorf("<< outside arith: got PLUSEQ, want LTLT")
		}
	}
	if !found {
		t.Errorf("<< outside arith: LTLT token not found in stream")
	}
}

// TestNonArithSlashPath verifies that / in a path outside arithmetic
// still lexes as part of an identifier (IDENT or SLASH), never PLUSEQ.
func TestNonArithSlashPath(t *testing.T) {
	toks := tokensFor(t, "echo /usr/bin/env\n")
	for _, tok := range toks {
		if tok.Type == token.PLUSEQ {
			t.Errorf("/ in path: got unexpected PLUSEQ token")
		}
	}
}

// TestNonArithAsteriskGlob verifies that * outside arithmetic remains
// ASTERISK (glob), not confused with **= handling.
func TestNonArithAsteriskGlob(t *testing.T) {
	toks := tokensFor(t, "echo *\n")
	found := false
	for _, tok := range toks {
		if tok.Type == token.ASTERISK {
			found = true
		}
		if tok.Type == token.PLUSEQ {
			t.Errorf("* outside arith: got unexpected PLUSEQ token")
		}
	}
	if !found {
		t.Errorf("* outside arith: ASTERISK token not found")
	}
}
