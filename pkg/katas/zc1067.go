package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1067",
		Title: "Separate `export` and assignment to avoid masking return codes",
		Description: "Running `export var=$(cmd)` masks the return code of `cmd`. " +
			"The exit status will be that of `export` (usually 0). " +
			"Declare the variable first or export it after assignment.",
		Check: checkZC1067,
	})
}

func checkZC1067(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if command is "export"
	name := cmd.Name.String()
	if name != "export" {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		if containsSubstitutionAfterEquals(arg) {
			violations = append(violations, Violation{
				KataID: "ZC1067",
				Message: "Exporting and assigning a command substitution in one step masks the return value. " +
					"Use `var=$(cmd); export var`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
			})
		}
	}

	return violations
}

func containsSubstitutionAfterEquals(expr ast.Expression) bool {
	// Debug: print node type and string representation
	// fmt.Printf("DEBUG: Type=%T String=%q\n", expr, expr.String())

	// Check if the argument contains an equals sign "="
	if stringIndex(expr.String(), "=") < 0 {
		return false
	}

	// Now check if it contains a command substitution
	return containsSubst(expr)
}

func stringIndex(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func isSubstitution(n ast.Node) bool {
	switch n.(type) {
	case *ast.CommandSubstitution, *ast.DollarParenExpression:
		return true
	}
	return false
}

func containsSubst(n ast.Node) bool {
	found := false
	ast.Walk(n, func(node ast.Node) bool {
		if isSubstitution(node) {
			found = true
			return false
		}
		return true
	})
	return found
}
