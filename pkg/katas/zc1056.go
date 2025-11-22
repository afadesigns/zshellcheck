package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1056",
		Title:       "Avoid `$((...))` as a statement",
		Description: "Using `$((...))` as a statement tries to execute the result as a command. Use `((...))` for arithmetic evaluation/assignment.",
		Check:       checkZC1056,
	})
}

func checkZC1056(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if the command Name is a DollarParenExpression (arithmetic)
	var dpe *ast.DollarParenExpression
	
	if d, ok := cmd.Name.(*ast.DollarParenExpression); ok {
		dpe = d
	} else if concat, ok := cmd.Name.(*ast.ConcatenatedExpression); ok {
		if len(concat.Parts) == 1 {
			if d, ok := concat.Parts[0].(*ast.DollarParenExpression); ok {
				dpe = d
			}
		}
	}

	if dpe == nil {
		return nil
	}

	// Check if it is an arithmetic expression, not a command substitution.
	// Our parser distinguishes:
	// $(( ... )) -> Command is usually Infix/Prefix/Identifier/Integer/Grouped
	// $( ... )   -> Command is usually SimpleCommand (via parseCommandList)
	
	isArithmetic := true
	
	switch dpe.Command.(type) {
	case *ast.SimpleCommand:
		// $(cmd)
		isArithmetic = false
	case *ast.ConcatenatedExpression:
		// $(cmd arg)
		isArithmetic = false
	}

	if isArithmetic {
		return []Violation{{
			KataID:  "ZC1056",
			Message: "Avoid `$((...))` as a statement. It executes the result. Use `((...))` for arithmetic.",
			Line:    dpe.TokenLiteralNode().Line,
			Column:  dpe.TokenLiteralNode().Column,
		}}
	}

	return nil
}
