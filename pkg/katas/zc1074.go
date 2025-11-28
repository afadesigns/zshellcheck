package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	kata := Kata{
		ID:    "ZC1074",
		Title: "Prefer modifiers :h/:t over dirname/basename",
		Description: "Zsh provides modifiers like `:h` (head/dirname) and `:t` (tail/basename) " +
			"that are faster and more idiomatic than spawning external commands.",
		Check: checkZC1074,
	}
	RegisterKata(ast.CommandSubstitutionNode, kata)
	RegisterKata(ast.DollarParenExpressionNode, kata)
}

func checkZC1074(node ast.Node) []Violation {
	var command ast.Node

	switch n := node.(type) {
	case *ast.CommandSubstitution:
		command = n.Command
	case *ast.DollarParenExpression:
		command = n.Command
	default:
		return nil
	}

	// Check if command is "dirname" or "basename"
	if cmd, ok := command.(*ast.SimpleCommand); ok {
		cmdName := cmd.Name.String()
		if cmdName == "dirname" {
			return []Violation{{
				KataID:  "ZC1074",
				Message: "Use '${var:h}' instead of '$(dirname $var)'. Modifiers are faster and built-in.",
				Line:    node.TokenLiteralNode().Line,
				Column:  node.TokenLiteralNode().Column,
			}}
		}
		if cmdName == "basename" {
			return []Violation{{
				KataID:  "ZC1074",
				Message: "Use '${var:t}' instead of '$(basename $var)'. Modifiers are faster and built-in.",
				Line:    node.TokenLiteralNode().Line,
				Column:  node.TokenLiteralNode().Column,
			}}
		}
	}

	return nil
}
