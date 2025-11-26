package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/token"
)

func init() {
	kata := Kata{
		ID:    "ZC1004",
		Title: "Use `return` instead of `exit` in functions",
		Description: "Using `exit` in a function terminates the entire shell, which is often unintended " +
			"in interactive sessions or sourced scripts. Use `return` to exit the function.",
		Check: checkZC1004,
	}
	RegisterKata(ast.FunctionDefinitionNode, kata)
	RegisterKata(ast.FunctionLiteralNode, kata)
}

func checkZC1004(node ast.Node) []Violation {
	var body ast.Statement

	switch n := node.(type) {
	case *ast.FunctionDefinition:
		body = n.Body
	case *ast.FunctionLiteral:
		body = n.Body
	default:
		return nil
	}

	violations := []Violation{}

	ast.Walk(body, func(n ast.Node) bool {
		// Stop traversal at subshell boundaries where exit is safe/scoped
		switch t := n.(type) {
		case *ast.GroupedExpression: // ( ... )
			return false
		case *ast.Subshell: // ( ... ) as subshell
			return false
		case *ast.CommandSubstitution: // ` ... `
			return false
		case *ast.DollarParenExpression: // $( ... )
			return false
		case *ast.BlockStatement:
			if t.Token.Type == token.LPAREN { // ( ... ) as a statement block
				return false
			}
		}

		if cmd, ok := n.(*ast.SimpleCommand); ok {
			if cmd.Name.String() == "exit" {
				violations = append(violations, Violation{
					KataID:  "ZC1004",
					Message: "Use `return` instead of `exit` in functions to avoid killing the shell.",
					Line:    cmd.TokenLiteralNode().Line,
					Column:  cmd.TokenLiteralNode().Column,
				})
			}
		}
		return true
	})

	return violations
}
