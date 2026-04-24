package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1003",
		Title: "Use `((...))` for arithmetic comparisons instead of `[` or `test`",
		Description: "Bash/Zsh have a dedicated arithmetic context `((...))` " +
			"which is cleaner and faster than `[` or `test` for numeric comparisons.",
		Severity: SeverityStyle,
		Check:    checkZC1003,
		Fix:      fixZC1003,
	})
}

// arithCmpReplacements maps the dash-prefixed test-comparison
// operators to their arithmetic equivalents. Used by fixZC1003 to
// rewrite `[ x -eq y ]` → `(( x == y ))`.
var arithCmpReplacements = map[string]string{
	"-eq": "==",
	"-ne": "!=",
	"-lt": "<",
	"-le": "<=",
	"-gt": ">",
	"-ge": ">=",
}

func checkZC1003(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		name := cmd.Name.String()
		if name == "[" || name == "test" {
			for _, arg := range cmd.Arguments {
				val := arg.String()
				// Trim parens added by AST String() method for expressions
				val = strings.Trim(val, "()")

				if _, found := arithCmpReplacements[val]; found {
					violations = append(violations, Violation{
						KataID:  "ZC1003",
						Message: "Use `((...))` for arithmetic comparisons instead of `[` or `test`.",
						Line:    cmd.Token.Line,
						Column:  cmd.Token.Column,
						Level:   SeverityStyle,
					})
					return violations
				}
			}
		}
	}

	return violations
}

// fixZC1003 rewrites `[ x -eq y ]` to `(( x == y ))`. Only the
// `[` form is auto-fixable: replace the opening `[` with `((`,
// the matching `]` with `))`, and the `-eq`/etc. operator with
// its arithmetic equivalent. The `test` form has no closing
// terminator on the line and the surrounding context (pipelines,
// chains) makes a safe rewrite ambiguous; it stays detection-only.
//
// Bails when the command shape doesn't match: more than one
// comparison operator, no `[` byte at the violation column,
// missing close bracket, or no recognised operator.
func fixZC1003(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok || cmd.Name == nil {
		return nil
	}
	if cmd.Name.String() != "[" {
		return nil
	}
	var op string
	var opIdx int
	for i, arg := range cmd.Arguments {
		val := strings.Trim(arg.String(), "()")
		if _, found := arithCmpReplacements[val]; found {
			if op != "" {
				return nil // multiple ops — not safe to auto-fix
			}
			op = val
			opIdx = i
		}
	}
	if op == "" {
		return nil
	}
	openOff := LineColToByteOffset(source, v.Line, v.Column)
	if openOff < 0 || openOff >= len(source) || source[openOff] != '[' {
		return nil
	}
	closeOff := findTestCloseBracket(source, openOff)
	if closeOff < 0 {
		return nil
	}
	opTok := cmd.Arguments[opIdx]
	opLine := opTok.TokenLiteralNode().Line
	opCol := opTok.TokenLiteralNode().Column
	return []FixEdit{
		{Line: v.Line, Column: v.Column, Length: 1, Replace: "(("},
		{Line: opLine, Column: opCol, Length: len(op), Replace: arithCmpReplacements[op]},
		offsetToEdit(source, closeOff, 1, "))"),
	}
}
