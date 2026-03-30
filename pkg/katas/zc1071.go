package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:          "ZC1071",
		Title:       "Use `+=` for appending to arrays",
		Description: "Appending to an array using `arr=($arr ...)` is verbose and slower. Use `arr+=(...)` instead.",
		Severity:    SeverityWarning,
		Check:       checkZC1071,
	})
}

func checkZC1071(node ast.Node) []Violation {
	infix, ok := node.(*ast.InfixExpression)
	if !ok || infix.Operator != "=" {
		return nil
	}

	ident, ok := infix.Left.(*ast.Identifier)
	if !ok {
		return nil
	}
	varName := ident.Value

	arrayLit, ok := infix.Right.(*ast.ArrayLiteral)
	if !ok {
		return nil
	}

	found := false
	checkNode := func(n ast.Node) bool {
		if found {
			return false
		}
		if aa, ok := n.(*ast.ArrayAccess); ok {
			if id, ok := aa.Left.(*ast.Identifier); ok && id.Value == varName {
				found = true
				return false
			}
		}
		if id, ok := n.(*ast.Identifier); ok {
			if id.Value == "$"+varName || id.Value == "${"+varName+"}" {
				found = true
				return false
			}
		}
		if prefix, ok := n.(*ast.PrefixExpression); ok {
			if prefix.Operator == "$" {
				if id, ok := prefix.Right.(*ast.Identifier); ok && id.Value == varName {
					found = true
					return false
				}
			}
		}
		return true
	}

	for _, elem := range arrayLit.Elements {
		if found {
			break
		}
		ast.Walk(elem, checkNode)
	}

	if found {
		leftToken := infix.Left.TokenLiteralNode()
		return []Violation{{
			KataID: "ZC1071",
			Message: "Appending to an array using `arr=($arr ...)` is verbose and slower. " +
				"Use `arr+=(...)` instead.",
			Line:   leftToken.Line,
			Column: leftToken.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
