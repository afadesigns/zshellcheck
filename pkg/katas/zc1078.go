package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1078",
		Title: "Quote `$@` and `$*` when passing arguments",
		Description: "Using unquoted `$@` or `$*` splits arguments by IFS (usually space). " +
			"Use `\"$@\"` to preserve the original argument grouping, or `\"$*\"` to join them into a single string.",
		Severity: SeverityWarning,
		Check:    checkZC1078,
		Fix:      fixZC1078,
	})
}

// fixZC1078 wraps an unquoted `$@` / `$*` argument in double-quotes.
// Both tokens are exactly two bytes; the two-edit insertion always
// surrounds the same 2-byte run.
func fixZC1078(_ ast.Node, v Violation, source []byte) []FixEdit {
	start := LineColToByteOffset(source, v.Line, v.Column)
	if start < 0 || start+2 > len(source) {
		return nil
	}
	if source[start] != '$' || (source[start+1] != '@' && source[start+1] != '*') {
		return nil
	}
	endLine, endCol := offsetLineColZC1078(source, start+2)
	if endLine < 0 {
		return nil
	}
	return []FixEdit{
		{Line: v.Line, Column: v.Column, Length: 0, Replace: `"`},
		{Line: endLine, Column: endCol, Length: 0, Replace: `"`},
	}
}

func offsetLineColZC1078(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1078(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		// Check string representation to catch various parsed forms of $@ and $*
		// unquoted $@ might be parsed as Identifier "$@" -> String() == "$@"
		// unquoted $* might be parsed as GroupedExpression -> String() == "($*)"
		// or other variations depending on parser state (e.g. PrefixExpression)

		s := arg.String()

		// Removing parens from GroupedExpression string representation for checking
		// (Note: String() adds parens for GroupedExpression)
		if len(s) >= 2 && s[0] == '(' && s[len(s)-1] == ')' {
			s = s[1 : len(s)-1]
		}

		if s == "$@" || s == "$*" {
			violations = append(violations, Violation{
				KataID:  "ZC1078",
				Message: "Unquoted " + s + " splits arguments. Use \"" + s + "\" to preserve structure.",
				Line:    arg.TokenLiteralNode().Line,
				Column:  arg.TokenLiteralNode().Column,
				Level:   SeverityWarning,
			})
		}
	}

	return violations
}
