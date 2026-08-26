// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

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

// HasArgFlag reports whether any cmd argument stringifies to a key
// present in flags. Replaces the recurring `for arg ... if v == "-x" ||
// v == "-y" || ...` chains that drive kata cyclomatic complexity above
// the gocyclo threshold.
func HasArgFlag(cmd *ast.SimpleCommand, flags map[string]struct{}) bool {
	for _, arg := range cmd.Arguments {
		if _, hit := flags[arg.String()]; hit {
			return true
		}
	}
	return false
}

// ArgValueAfter returns the stringified argument that immediately
// follows the first argument matching key, or the empty string when
// the key is absent or sits at the tail. Used by katas that match
// `--bind 0.0.0.0` style options.
func ArgValueAfter(cmd *ast.SimpleCommand, key string) string {
	for i, arg := range cmd.Arguments {
		if arg.String() == key && i+1 < len(cmd.Arguments) {
			return cmd.Arguments[i+1].String()
		}
	}
	return ""
}

// CommandIdentifier returns the head identifier value of cmd, or "" if
// the head is not an identifier. Wraps the common type-assertion guard
// at the top of every Check function.
func CommandIdentifier(cmd *ast.SimpleCommand) string {
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return ""
	}
	return ident.Value
}

// LongFlagName splits a long option into its name and reports whether the
// argument carried an inline value. `--depth=1` yields `--depth`, true;
// `--depth` yields `--depth`, false. A word that is not a long option comes
// back unchanged. Most GNU-style tools accept both spellings for a
// value-taking option, so a kata that compares an argument against a bare
// `--flag` sees only half of real usage.
func LongFlagName(word string) (string, bool) {
	if !strings.HasPrefix(word, "--") {
		return word, false
	}
	if name, _, found := strings.Cut(word, "="); found {
		return name, true
	}
	return word, false
}

// FlagValueAt returns the value carried by the option at args[i] when that
// option is one of names, together with the number of arguments it spans.
// `--directory=/` spans one argument, `--directory /` spans two. A word that
// names no listed option, or an option whose value never arrives, yields
// ("", 0).
func FlagValueAt(args []ast.Expression, i int, names ...string) (string, int) {
	if i < 0 || i >= len(args) {
		return "", 0
	}
	word := args[i].String()
	name, inline := LongFlagName(word)
	matched := false
	for _, want := range names {
		if name == want {
			matched = true
			break
		}
	}
	if !matched {
		return "", 0
	}
	if inline {
		_, value, _ := strings.Cut(word, "=")
		return value, 1
	}
	if i+1 >= len(args) {
		return "", 0
	}
	return args[i+1].String(), 2
}
