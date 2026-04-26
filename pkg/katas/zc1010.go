// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1010",
		Title: "Use [[ ... ]] instead of [ ... ]",
		Description: "Zsh's [[ ... ]] is more powerful and safer than [ ... ]. " +
			"It supports pattern matching, regex, and doesn't require quoting variables to prevent word splitting.",
		Severity: SeverityStyle,
		Check:    checkZC1010,
		Fix:      fixZC1010,
	})
}

// fixZC1010 rewrites a `[ … ]` test command to `[[ … ]]`. The opening
// bracket at the violation's coordinates becomes `[[`; the matching
// closing bracket on the same logical line becomes `]]`. Contents
// stay byte-identical so quoting and expansions are preserved.
//
// Bail when the shape is not a simple `[ … ]` test (e.g. second token
// is not `[`, or the logical line has no closing bracket): a
// malformed test is not safely auto-fixable.
func fixZC1010(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if cmd.Name == nil || cmd.Name.String() != "[" {
		return nil
	}
	open := LineColToByteOffset(source, v.Line, v.Column)
	if open < 0 || open >= len(source) || source[open] != '[' {
		return nil
	}
	close := findTestCloseBracket(source, open)
	if close < 0 {
		return nil
	}
	return []FixEdit{
		{Line: v.Line, Column: v.Column, Length: 1, Replace: "[["},
		offsetToEdit(source, close, 1, "]]"),
	}
}

// findTestCloseBracket returns the byte offset of the closing `]`
// that terminates a `[ … ]` test opened at open, or -1 when no clean
// close is found before an end-of-statement terminator. Single- and
// double-quoted strings, and `${…}` braces are respected so the scan
// doesn't match `]` inside `"$arr[idx]"`.
func findTestCloseBracket(source []byte, open int) int {
	st := bracketScan{}
	for i := open + 1; i < len(source); i++ {
		if next, ok := st.advance(source, i); ok {
			i = next
			continue
		}
		if st.closedAt(source, i) {
			return i
		}
		if st.terminated(source, i) {
			return -1
		}
	}
	return -1
}

type bracketScan struct {
	inSingle, inDouble bool
	braceDepth         int
}

// advance returns (newIndex, true) when the byte at i drives the scanner
// forward past quoted/escaped material without further classification.
func (s *bracketScan) advance(source []byte, i int) (int, bool) {
	c := source[i]
	switch {
	case s.inSingle:
		if c == '\'' {
			s.inSingle = false
		}
		return i, true
	case s.inDouble:
		if c == '\\' && i+1 < len(source) {
			return i + 1, true
		}
		if c == '"' {
			s.inDouble = false
		}
		return i, true
	}
	return i, false
}

// closedAt reports whether the unquoted byte at i is the matching `]`.
// Side effect: classifies opening / closing braces to track depth.
func (s *bracketScan) closedAt(source []byte, i int) bool {
	switch source[i] {
	case '\'':
		s.inSingle = true
	case '"':
		s.inDouble = true
	case '{':
		s.braceDepth++
	case '}':
		if s.braceDepth > 0 {
			s.braceDepth--
		}
	case ']':
		if s.braceDepth == 0 {
			return true
		}
	}
	return false
}

func (s *bracketScan) terminated(source []byte, i int) bool {
	if s.inSingle || s.inDouble {
		return false
	}
	c := source[i]
	return c == '\n' || c == ';'
}

// offsetToEdit builds a FixEdit whose Line/Column correspond to the
// given byte offset inside source. Used when a Fix already has a
// byte offset but the FixEdit type expects 1-based coordinates.
func offsetToEdit(source []byte, offset, length int, replace string) FixEdit {
	line := 1
	col := 1
	for i := 0; i < offset && i < len(source); i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return FixEdit{Line: line, Column: col, Length: length, Replace: replace}
}

func checkZC1010(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		// Check if command name is "["
		if cmd.Name.String() == "[" {
			violations = append(violations, Violation{
				KataID:  "ZC1010",
				Message: "Use `[[ ... ]]` instead of `[ ... ]` or `test`. `[[` is safer and more powerful.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityStyle,
			})
		}
	}

	return violations
}
