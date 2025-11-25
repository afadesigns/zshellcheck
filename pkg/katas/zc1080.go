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
		// Check if item looks like a glob pattern (contains *, ?, [])
		// and DOES NOT contain (N) or (nullglob) qualifier.
		
		// Parser might return it as Identifier (if simple glob) or Prefix/Infix (if paths).
		// Or just String().
		
		s := item.String()
		
		// Simple heuristic for glob characters
		isGlob := strings.ContainsAny(s, "*?[]")
		
		if isGlob {
			// Check for existing nullglob qualifier
			// (N) is standard short form.
			// We check if it ends with (N) or contains (N) in qualifiers.
			// Zsh qualifiers are usually at the end in parens.
			
			if !strings.Contains(s, "(N)") && !strings.Contains(s, "N") { 
				// "N" checking is tricky because (N) is the syntax. 
				// If we have qualifiers like (*.txt)(.N), we check for N inside parens.
				// But parsing qualifiers without a full zsh lexer is hard.
				// Let's stick to checking for explicit "(N)" or "N" inside the last parenthesized group?
				// Or simplistically: if it doesn't contain "(N)", warn.
				
				// Refined check: 
				// Must contain glob char.
				// Must NOT contain "(N)".
				
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
