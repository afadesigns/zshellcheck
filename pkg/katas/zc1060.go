package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1060",
		Title: "Avoid `ps | grep` without exclusion",
		Description: "`ps | grep pattern` often matches the grep process itself. " +
			"Use `grep [p]attern`, `pgrep`, or exclude grep with `grep -v grep`.",
		Check: checkZC1060,
	})
}

func checkZC1060(node ast.Node) []Violation {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}

	// Check if left command is `ps`
	if !isCommandName(pipe.Left, "ps") {
		return nil
	}

	// Check if right command is `grep`
	if !isCommandName(pipe.Right, "grep") {
		return nil
	}

	// Check if grep arguments exclude the grep process
	// Strategies:
	// 1. `grep -v grep` (chained pipe?)
	//    If pipe.Right is `grep`, we only see `grep ...`.
	//    If user does `ps | grep foo | grep -v grep`, the parsing structure is `(ps | grep foo) | grep -v grep`.
	//    So we are looking at `ps | grep foo`. The parent pipe handles the exclusion?
	//    We can't see the parent here easily.
	//    BUT, `ps | grep foo` is inherently risky unless `foo` uses `[]`.
	// 2. Pattern uses `[]`. e.g. `grep [f]oo`.

	// We inspect `grep` arguments.
	cmd, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok {
		return nil // complex command
	}

	hasExclusion := false

	for _, arg := range cmd.Arguments {
		// Check if arg contains `[` and `]`.
		val := getStringValueZC1060(arg)
		// Naive check for `[...]` pattern
		// If it starts with - (flag), ignore unless it is -v grep?
		// But we only check THIS grep.

		if len(val) > 0 && val[0] != '-' {
			// Assume this is the pattern
			// Check for brackets
			for i := 0; i < len(val); i++ {
				if val[i] == '[' {
					// Look for closing ]
					for j := i + 1; j < len(val); j++ {
						if val[j] == ']' {
							hasExclusion = true
							break
						}
					}
				}
			}
		}
	}

	if !hasExclusion {
		return []Violation{{
			KataID:  "ZC1060",
			Message: "`ps | grep pattern` matches the grep process itself. Use `grep [p]attern` to exclude the grep process.",
			Line:    pipe.TokenLiteralNode().Line,
			Column:  pipe.TokenLiteralNode().Column,
		}}
	}

	return nil
}

func isCommandName(node ast.Node, name string) bool {
	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if ident, ok := cmd.Name.(*ast.Identifier); ok {
			return ident.Value == name
		}
	}
	return false
}

func getStringValueZC1060(node ast.Node) string {
	switch n := node.(type) {
	case *ast.StringLiteral:
		return n.Value
	case *ast.ConcatenatedExpression:
		// Simplify
		return ""
	case *ast.Identifier:
		return n.Value
	}
	return ""
}
