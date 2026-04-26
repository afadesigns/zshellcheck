// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import "github.com/afadesigns/zshellcheck/pkg/ast"

// LineColToByteOffset converts a 1-based (line, column) coordinate
// pair into a 0-based byte offset within source. It returns -1 when
// the coordinates are out of range. Used by kata Fix functions that
// need to splice source around a token known by its line/column only.
func LineColToByteOffset(source []byte, line, col int) int {
	if line < 1 || col < 1 {
		return -1
	}
	curLine := 1
	curCol := 1
	for i, b := range source {
		if curLine == line && curCol == col {
			return i
		}
		if b == '\n' {
			curLine++
			curCol = 1
			continue
		}
		curCol++
	}
	// End-of-file case: coordinates pointing at the byte just past the
	// last character are valid when the last line has no newline.
	if curLine == line && curCol == col {
		return len(source)
	}
	return -1
}

// IdentLenAt returns the length in bytes of the identifier starting at
// source[offset]. An identifier is a run of [A-Za-z0-9_-]. Returns 0
// when offset is out of range or does not start on an identifier byte.
// Useful when a kata wants to replace the command name at the head of
// a SimpleCommand (Token coordinates point at the name start).
func IdentLenAt(source []byte, offset int) int {
	if offset < 0 || offset >= len(source) {
		return 0
	}
	n := 0
	for offset+n < len(source) && isIdentByte(source[offset+n]) {
		n++
	}
	return n
}

func isIdentByte(b byte) bool {
	switch {
	case b >= 'a' && b <= 'z':
		return true
	case b >= 'A' && b <= 'Z':
		return true
	case b >= '0' && b <= '9':
		return true
	case b == '_' || b == '-':
		return true
	}
	return false
}

// FlagArgPosition scans cmd.Arguments for the first arg whose
// String() value matches one of needles. Returns the (line, column)
// of that arg's leading token. When no arg matches, falls back to
// the cmd.Token coordinates so callers always have something to
// report. Used by katas that detect dangerous long-flags
// (`--delete-secret-keys`, `--bind_ip 0.0.0.0`, etc.) and want the
// violation pointer at the flag itself, not the host command name.
func FlagArgPosition(cmd *ast.SimpleCommand, needles map[string]bool) (int, int) {
	for _, arg := range cmd.Arguments {
		if needles[arg.String()] {
			tok := arg.TokenLiteralNode()
			return tok.Line, tok.Column
		}
	}
	return cmd.Token.Line, cmd.Token.Column
}
