package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.ForLoopStatementNode, Kata{
		ID:    "ZC1080",
		Title: "Use `(N)` nullglob qualifier for globs in loops",
		Description: "In Zsh, if a glob matches no files, it throws an error by default. " +
			"When iterating over a glob in a `for` loop, use the `(N)` glob qualifier to allow it to match nothing (nullglob).",
		Check: checkZC1080,
	})
}

func checkZC1080(node ast.Node) []Violation {
	loop, ok := node.(*ast.ForLoopStatement)
	if !ok {
		return nil
	}

	// Only check for-each loops: for i in items...
	if loop.Items == nil {
		return nil // C-style loop
	}

	violations := []Violation{}

	for _, item := range loop.Items {
		if hasGlobChars(item) {
			s := item.String()
			if !strings.Contains(s, "(N)") && !strings.Contains(s, "N") {
				violations = append(violations, Violation{
					KataID:  "ZC1080",
					Message: "Glob '" + s + "' will error if no matches found. Append `(N)` to make it nullglob.",
					Line:    item.TokenLiteralNode().Line,
					Column:  item.TokenLiteralNode().Column,
				})
			}
		}
	}

	return violations
}

func hasGlobChars(node ast.Node) bool {
	// Check AST nodes to see if they contain exposed glob characters (*, ?, [)
	switch n := node.(type) {
	case *ast.StringLiteral:
		// Check for unquoted glob chars in the literal
		// If parsed as StringLiteral, it might be single quoted (no glob) or simple word (glob).
		// Lexer strips quotes from Literal? No, value usually keeps them.
		val := n.Value
		if len(val) >= 2 && (val[0] == '\'' || val[0] == '"') {
			return false // Quoted strings don't glob
		}
		return strings.ContainsAny(val, "*?[]")
	case *ast.Identifier:
		// Identifiers don't glob unless they contain * (which usually makes them NOT identifiers but string/prefix)
		// But Parser might be lenient.
		return strings.ContainsAny(n.Value, "*?[]")
	case *ast.PrefixExpression:
		// *, ? prefix operators
		if n.Operator == "*" || n.Operator == "?" {
			return true
		}
		// Recursive check
		return hasGlobChars(n.Right)
	case *ast.ConcatenatedExpression:
		for _, part := range n.Parts {
			if hasGlobChars(part) {
				return true
			}
		}
		return false
	case *ast.ArrayAccess:
		return false // Array access ${...} is not a file glob
	case *ast.SimpleCommand:
		// [ char range ] is parsed as SimpleCommand sometimes?
		// No, usually Concatenated or StringLiteral if [ is treated as literal.
		// If [ is SimpleCommand name (e.g. `[` test command), it's not a glob.
		if n.Name.String() == "[" {
			// This is test command or char range start?
			// parseCommandWord treats [ as literal string if in exception list.
			// So it won't be SimpleCommand here.
		}
	}
	return false
}
