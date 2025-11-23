package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.FunctionDefinitionNode, Kata{
		ID:          "ZC1043",
		Title:       "Use `local` for variables in functions",
		Description: "Variables defined in functions are global by default in Zsh. " +
			"Use `local` to scope them to the function.",
		Check: checkZC1043,
	})
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
						violations = append(violations, Violation{
							KataID:  "ZC1043",
							Message: "Variable '" + ident.Value + "' is assigned without 'local'. It will be global. " +
								"Use `local " + ident.Value + "=" + assign.Right.String() + "`.",
							Line:   ident.Token.Line,
							Column: ident.Token.Column,
						})
					}
				}
			}
		}

		return true
	})

	return violations
}
