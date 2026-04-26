// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.FunctionDefinitionNode, Kata{
		ID:    "ZC1043",
		Title: "Use `local` for variables in functions",
		Description: "Variables defined in functions are global by default in Zsh. " +
			"Use `local` to scope them to the function.",
		Severity: SeverityStyle,
		Check:    checkZC1043,
		Fix:      fixZC1043,
	})
}

// fixZC1043 prepends `local ` to the unscoped assignment the detector
// flagged. The violation's Line/Column points at the assignment LHS;
// inserting `local ` there yields `local NAME=value`. On re-run the
// detector recognises `local …` as a declaration and skips the line,
// so the rewrite is idempotent.
func fixZC1043(_ ast.Node, v Violation, source []byte) []FixEdit {
	off := LineColToByteOffset(source, v.Line, v.Column)
	if off < 0 || off >= len(source) {
		return nil
	}
	// Defensive idempotency guard: refuse to insert if a declaration
	// keyword already sits at the violation column.
	for _, prefix := range []string{"local ", "typeset ", "declare ", "integer ", "float ", "readonly "} {
		end := off + len(prefix)
		if end <= len(source) && string(source[off:end]) == prefix {
			return nil
		}
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  0,
		Replace: "local ",
	}}
}

func checkZC1043(node ast.Node) []Violation {
	funcDef, ok := node.(*ast.FunctionDefinition)
	if !ok {
		return nil
	}

	violations := []Violation{}
	locals := make(map[string]bool)

	ast.Walk(funcDef.Body, func(n ast.Node) bool {
		// Stop walking into nested function definitions
		if _, ok := n.(*ast.FunctionDefinition); ok && n != funcDef {
			return false
		}

		// Track local declarations
		if cmd, ok := n.(*ast.SimpleCommand); ok {
			nameStr := cmd.Name.String()
			if nameStr == "local" || nameStr == "typeset" || nameStr == "declare" ||
				nameStr == "integer" || nameStr == "float" || nameStr == "readonly" {
				for _, arg := range cmd.Arguments {
					// Arg can be "x" or "x=1" or "-r"
					argStr := arg.String()
					if len(argStr) > 0 && argStr[0] == '-' {
						continue // Skip options
					}
					// Extract name before '='
					varName := argStr
					for i, c := range argStr {
						if c == '=' {
							varName = argStr[:i]
							break
						}
					}
					locals[varName] = true
				}
			}
		}

		// Check assignments
		if exprStmt, ok := n.(*ast.ExpressionStatement); ok {
			if assign, ok := exprStmt.Expression.(*ast.InfixExpression); ok && assign.Operator == "=" {
				if ident, ok := assign.Left.(*ast.Identifier); ok {
					if !locals[ident.Value] {
						// Empty RHS (`VAR=` at end of line) is valid Zsh
						// and the parser records it with Right == nil.
						// Fall back to an empty string so the message
						// builder doesn't deref nil.
						rhs := ""
						if assign.Right != nil {
							rhs = assign.Right.String()
						}
						violations = append(violations, Violation{
							KataID: "ZC1043",
							Message: "Variable '" + ident.Value + "' is assigned without 'local'. It will be global. " +
								"Use `local " + ident.Value + "=" + rhs + "`.",
							Line:   ident.Token.Line,
							Column: ident.Token.Column,
							Level:  SeverityStyle,
						})
					}
				}
			}
		}

		return true
	})

	return violations
}
