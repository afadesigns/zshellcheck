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
	inSingle := false
	inDouble := false
	braceDepth := 0
	for i := open + 1; i < len(source); i++ {
		c := source[i]
		switch {
		case inSingle:
			if c == '\'' {
				inSingle = false
			}
		case inDouble:
			if c == '\\' && i+1 < len(source) {
				i++
				continue
			}
			if c == '"' {
				inDouble = false
			}
		default:
			switch c {
			case '\'':
				inSingle = true
			case '"':
				inDouble = true
			case '{':
				braceDepth++
			case '}':
				if braceDepth > 0 {
					braceDepth--
				}
			case '\n', ';':
				return -1
			case ']':
				if braceDepth == 0 {
					return i
				}
			}
		}
	}
	return -1
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
