package katas

import (
	// "fmt"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.ExpressionStatementNode, Kata{
		ID:          "ZC1071",
		Title:       "Use `+=` for appending to arrays",
		Description: "Appending to an array using `arr=($arr ...)` is verbose and slower. Use `arr+=(...)` instead.",
		Check:       checkZC1071,
	})
}

func checkZC1071(node ast.Node) []Violation {
	exprStmt, ok := node.(*ast.ExpressionStatement)
	if !ok {
		return nil
	}

	infixExpr, ok := exprStmt.Expression.(*ast.InfixExpression)
	if !ok || infixExpr.Operator != "=" {
		return nil
	}

	leftIdent, ok := infixExpr.Left.(*ast.Identifier)
	if !ok {
		return nil
	}

	varName := leftIdent.Value
	valueExpr := infixExpr.Right

	if checkSelfReference(varName, valueExpr) {
		return []Violation{{
			KataID:  "ZC1071",
			Message: "Appending to an array using `arr=($arr ...)` is verbose and slower. Use `arr+=(...)` instead.",
			Line:    exprStmt.Token.Line,
			Column:  exprStmt.Token.Column,
		}}
	}

	return nil
}

func checkSelfReference(varName string, expr ast.Expression) bool {
	found := false
	checkNode := func(n ast.Node) bool {
		if ident, ok := n.(*ast.Identifier); ok {
			if ident.Value == varName || (strings.HasPrefix(ident.Value, "$") && strings.TrimPrefix(ident.Value, "$") == varName) {
				found = true
				return false
			}
		}
		if aa, ok := n.(*ast.ArrayAccess); ok {
			if ident, ok := aa.Left.(*ast.Identifier); ok && ident.Value == varName {
				found = true				
				return false
			}
		}
			if prefix, ok := n.(*ast.PrefixExpression); ok && prefix.Operator == "$" {
				if ident, ok := prefix.Right.(*ast.Identifier); ok && ident.Value == varName {
					found = true
					return false
				}
			}
		}
		return true
	}
	ast.Walk(expr, checkNode)
	return found
}
