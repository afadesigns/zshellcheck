// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
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
		Severity: SeverityInfo,
		Check:    checkZC1045,
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
						Level:  SeverityInfo,
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
					Level:  SeverityInfo,
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
	case *ast.CommandSubstitution, *ast.DollarParenExpression:
		return true
	case *ast.ConcatenatedExpression:
		return zc1045ConcatHasSub(n)
	case *ast.StringLiteral:
		return zc1045StringHasSub(n.Value)
	}
	return false
}

func zc1045ConcatHasSub(n *ast.ConcatenatedExpression) bool {
	for _, p := range n.Parts {
		if isCommandSubstitution(p) {
			return true
		}
	}
	return false
}

// zc1045StringHasSub scans a double-quoted string literal for embedded
// `$(...)` or backtick command substitutions, ignoring backslash-
// escaped bytes.
func zc1045StringHasSub(val string) bool {
	if len(val) < 2 || val[0] != '"' || val[len(val)-1] != '"' {
		return false
	}
	for i := 0; i < len(val); i++ {
		switch val[i] {
		case '\\':
			i++
		case '`':
			return true
		case '$':
			if i+1 < len(val) && val[i+1] == '(' {
				return true
			}
		}
	}
	return false
}
