package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1038",
		Title: "Avoid useless use of cat",
		Description: "Using `cat file | command` is unnecessary and inefficient. " +
			"Most commands can read from a file directly, e.g., `command file`. " +
			"If not, you can use input redirection: `command < file`.",
		Check: checkZC1038,
	})
}

func checkZC1038(node ast.Node) []Violation {
	violations := []Violation{}

	infix, ok := node.(*ast.InfixExpression)
	if !ok {
		return violations
	}
	
	if infix.Operator != "|" {
		return violations
	}

	cmd, ok := infix.Left.(*ast.SimpleCommand)
	if !ok {
		return violations
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return violations
	}
	
	if ident.Value != "cat" {
		return violations
	}

	// cat must have exactly one argument to be considered a "useless use" in this context.
	// cat without args reads from stdin (valid pipe).
	// cat with multiple args concatenates (valid use).
	if len(cmd.Arguments) == 1 {
		violations = append(violations, Violation{
			KataID: "ZC1038",
			Message: "Avoid useless use of cat. " +
				"Prefer `command file` or `command < file` over `cat file | command`.",
			Line:   ident.Token.Line,
			Column: ident.Token.Column,
		})
	}

	return violations
}
