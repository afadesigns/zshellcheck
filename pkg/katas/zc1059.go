package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1059",
		Title:       "Use `${var:?}` for `rm` arguments",
		Description: "Deleting a directory based on a variable is dangerous if the variable is empty or unset. Use `${var:?}` to fail if empty, or check explicitly.",
		Check:       checkZC1059,
	})
}

func checkZC1059(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if name, ok := cmd.Name.(*ast.Identifier); !ok || name.Value != "rm" {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		isUnsafeVar := false
		
		switch n := arg.(type) {
		case *ast.PrefixExpression:
			if n.Operator == "$" {
				isUnsafeVar = true // $VAR
			}
		case *ast.ArrayAccess:
			// ${VAR}. Check if it has modifiers?
			// Parser for ArrayAccess currently parses ${VAR} as ArrayAccess with Index=nil.
			// It does NOT parse modifiers like :?.
			// If the source has ${VAR:?}, the parser might fail or parse differently?
			// Current parser implementation for ArrayAccess:
			// Expects IDENT. Then optional [index]. Then }.
			// It does NOT handle : modifiers.
			// So ${VAR:?} would likely fail parsing or be parsed incorrectly.
			// If parser fails, we can't check it.
			// Assuming parser parses simple ${VAR}, we flag it.
			isUnsafeVar = true
		case *ast.StringLiteral:
			// "$VAR". 
			// If value is exactly "$VAR" or "${VAR}".
			// If value contains other things, it's safer (e.g. "$VAR/foo").
			// But "$VAR/" is dangerous too if VAR is empty.
			// For now, focus on exact variable.
			if isSimpleVariableString(n.Value) {
				isUnsafeVar = true
			}
		}

		if isUnsafeVar {
			violations = append(violations, Violation{
				KataID:  "ZC1059",
				Message: "Use `${var:?}` or ensure the variable is set before using it in `rm`.",
				Line:    arg.TokenLiteralNode().Line,
				Column:  arg.TokenLiteralNode().Column,
			})
		}
	}

	return violations
}

func isSimpleVariableString(s string) bool {
	// Check if string is "$VAR" or "${VAR}" (quoted)
	// Quotes are included in StringLiteral value.
	// "$VAR" -> len >= 4. e.g. "$V"
	if len(s) < 4 {
		return false
	}
	if s[0] != '"' || s[len(s)-1] != '"' {
		return false
	}
	inner := s[1 : len(s)-1]
	if len(inner) < 2 || inner[0] != '$' {
		return false
	}
	// Check if rest is valid identifier char (naive)
	// OR ${...}
	if inner[1] == '{' {
		// Must end with }
		if inner[len(inner)-1] != '}' {
			return false
		}
		return true // Assume ${...}
	}
	// $VAR
	return true
}
