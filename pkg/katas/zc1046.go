package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1046",
		Title: "Avoid `eval`",
		Description: "`eval` is dangerous as it executes arbitrary code. " +
			"Use arrays, parameter expansion, or other constructs instead.",
		Check: checkZC1046,
	})
}

func checkZC1046(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	name := cmd.Name.String()

	// Check for direct 'eval'
	if name == "eval" {
		return []Violation{{
			KataID:  "ZC1046",
			Message: "Avoid `eval`. It allows execution of arbitrary code and is hard to debug.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
		}}
	}

	// Check for 'builtin eval' or 'command eval'
	if (name == "builtin" || name == "command") && len(cmd.Arguments) > 0 {
		arg := cmd.Arguments[0]
		if arg.String() == "eval" {
			return []Violation{{
				KataID:  "ZC1046",
				Message: "Avoid `eval`. It allows execution of arbitrary code and is hard to debug.",
				Line:    arg.TokenLiteralNode().Line,
				Column:  arg.TokenLiteralNode().Column,
			}}
		}
	}

	return nil
}
