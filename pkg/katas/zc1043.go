package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.FunctionDefinitionNode, Kata{
		ID:          "ZC1043",
		Title:       "Use `local` for variables in functions",
		Description: "Variables defined in functions are global by default in Zsh. Use `local` to scope them to the function.",
		Check:       checkZC1043,
	})
}

func checkZC1043(node ast.Node) []Violation {
	funcDef, ok := node.(*ast.FunctionDefinition)
	if !ok {
		return nil
	}

	violations := []Violation{}

	ast.Walk(funcDef.Body, func(n ast.Node) bool {
		// We only care about "naked" assignments, which are ExpressionStatements containing an assignment.
		// Assignments inside `local x=1` are part of SimpleCommand and handled differently (not as naked ExpressionStatement).
		
		if exprStmt, ok := n.(*ast.ExpressionStatement); ok {
			if assign, ok := exprStmt.Expression.(*ast.InfixExpression); ok && assign.Operator == "=" {
				if ident, ok := assign.Left.(*ast.Identifier); ok {
					violations = append(violations, Violation{
						KataID:  "ZC1043",
						Message: "Variable '" + ident.Value + "' is assigned without 'local'. It will be global. Use `local " + ident.Value + "=" + assign.Right.String() + "`.",
						Line:    ident.Token.Line,
						Column:  ident.Token.Column,
					})
				}
			}
		}
		
		// Stop walking into nested function definitions
		if _, ok := n.(*ast.FunctionDefinition); ok && n != funcDef {
			return false
		}
		
		return true
	})
		if _, ok := n.(*ast.FunctionDefinition); ok && n != funcDef {
			return false
		}
		
		return true
	})

	return violations
}