package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	kata := Kata{
		ID:    "ZC1045",
		Title: "Declare and assign separately to avoid masking return values",
		Description: "Declaring a variable with `local var=$(cmd)` masks the return value of `cmd`. " +
			"The `local` command returns 0 (success) even if `cmd` fails. " +
			"Declare the variable first (`local var`), then assign it (`var=$(cmd)`).",
		Check: checkZC1045,
	}
	RegisterKata(ast.SimpleCommandNode, kata)
	RegisterKata(ast.DeclarationStatementNode, kata)
}

func checkZC1045(node ast.Node) []Violation {
	violations := []Violation{}

	// Check SimpleCommand (local, readonly)
	if cmd, ok := node.(*ast.SimpleCommand); ok {
		name := cmd.Name.String()
		if name == "local" || name == "readonly" {
			for _, arg := range cmd.Arguments {
				if hasCommandSubstitutionAssignment(arg) {
					violations = append(violations, Violation{
						KataID: "ZC1045",
						Message: "Declare and assign separately to avoid masking return values. " +
							"`" + name + " var=$(cmd)` masks the exit code of `cmd`.",
						Line:   arg.TokenLiteralNode().Line,
						Column: arg.TokenLiteralNode().Column,
					})
				}
			}
		}
	}

	// Check DeclarationStatement (typeset, declare)
	if decl, ok := node.(*ast.DeclarationStatement); ok {
		// Command is "typeset" or "declare"
		for _, assign := range decl.Assignments {
			if assign.Value != nil && isCommandSubstitution(assign.Value) {
				violations = append(violations, Violation{
					KataID: "ZC1045",
					Message: "Declare and assign separately to avoid masking return values. " +
						"`" + decl.Command + " var=$(cmd)` masks the exit code of `cmd`.",
					Line:   decl.Token.Line,
					Column: decl.Token.Column,
				})
			}
		}
	}

	return violations
}

func hasCommandSubstitutionAssignment(arg ast.Expression) bool {
	// Argument structure depends on parsing.
	// Usually ConcatenatedExpression for `var=$(cmd)`: [Identifier(var), StringLiteral(=), DollarParenExpression]
	// Or `var=`cmd``: [Identifier(var), StringLiteral(=), CommandSubstitution]

	concat, ok := arg.(*ast.ConcatenatedExpression)
	if !ok {
		return false
	}

	hasEquals := false
	hasCmdSubst := false

	for _, part := range concat.Parts {
		if str, ok := part.(*ast.StringLiteral); ok && str.Value == "=" {
			hasEquals = true
			continue
		}

		if hasEquals {
			// Check if RHS has command substitution
			if isCommandSubstitution(part) {
				hasCmdSubst = true
			}
		}
	}

	return hasEquals && hasCmdSubst
}

func isCommandSubstitution(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.CommandSubstitution:
		return true
	case *ast.DollarParenExpression:
		return true
	case *ast.ConcatenatedExpression:
		// Recursively check parts? e.g. `var="foo $(cmd)"`
		for _, p := range n.Parts {
			if isCommandSubstitution(p) {
				return true
			}
		}
	case *ast.StringLiteral:
		// Check for interpolation in double-quoted strings
		val := n.Value
		if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
			// Scan for $(...) or `...`
			// Simple heuristic: unescaped $ followed by ( or unescaped `
			for i := 0; i < len(val); i++ {
				if val[i] == '\\' {
					i++ // skip next
					continue
				}
				if val[i] == '`' {
					return true
				}
				if val[i] == '$' && i+1 < len(val) && val[i+1] == '(' {
					return true
				}
			}
		}
		return false
	}
	return false
}
