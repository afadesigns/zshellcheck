package katas

import (
	"regexp"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1054",
		Title:       "Use POSIX classes in regex/glob",
		Description: "Ranges like `[a-z]` are locale-dependent. Use `[[:lower:]]` or `[a-z]` with `LC_ALL=C` to be explicit.",
		Check:       checkZC1054,
	})
}

var rangeRegex = regexp.MustCompile(`\[[a-zA-Z0-9]-[a-zA-Z0-9]\]`)

func checkZC1054(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		val := getStringValueZC1054(arg)
		if rangeRegex.MatchString(val) {
			// Avoid flagging if it looks like a POSIX class like [[:lower:]]
			// But regex `\[[a-z]-[a-z]\]` matches `[a-z]` but not `[[:lower:]]`
			// Wait, `[[:lower:]]` contains `[:` which is not `[a-z]-[a-z]`.
			// So it should be safe.

			violations = append(violations, Violation{
				KataID:  "ZC1054",
				Message: "Ranges like `[a-z]` are locale-dependent. Use POSIX classes like `[[:lower:]]` or `[[:digit:]]`.",
				Line:    arg.TokenLiteralNode().Line,
				Column:  arg.TokenLiteralNode().Column,
			})
		}
	}

	return violations
}

func getStringValueZC1054(node ast.Node) string {
	switch n := node.(type) {
	case *ast.StringLiteral:
		return n.Value
	case *ast.Identifier:
		return n.Value
	case *ast.PrefixExpression:
		// Reconstruct -z, !foo, etc.
		// Right is Expression.
		return n.Operator + getStringValueZC1054(n.Right)
	case *ast.ConcatenatedExpression:
		var sb strings.Builder
		for _, p := range n.Parts {
			sb.WriteString(getStringValueZC1054(p))
		}
		return sb.String()
	}
	return node.String() // Fallback to String() for other types
}
