package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1089",
		Title: "Redirection order matters (`2>&1 > file`)",
		Description: "Redirecting stderr to stdout (`2>&1`) before redirecting stdout to a file (`> file`) " +
			"means stderr goes to the *original* stdout (usually tty), not the file. " +
			"Use `> file 2>&1` or `&> file` to redirect both.",
		Check: checkZC1089,
	})
}

func checkZC1089(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	idx2to1 := -1
	idxRedirect := -1
	var redirectArg ast.Expression

	for i, arg := range cmd.Arguments {
		s := arg.String()
		if s == "2>&1" {
			if idx2to1 == -1 {
				idx2to1 = i
			}
		} else if s == ">" || s == ">>" {
			// Found redirection operator
			if idxRedirect == -1 {
				idxRedirect = i
				redirectArg = arg
			}
		}
	}

	if idx2to1 != -1 && idxRedirect != -1 && idx2to1 < idxRedirect {
		return []Violation{
			{
				KataID:  "ZC1089",
				Message: "Redirection order matters. `2>&1 > file` does not redirect stderr to file. Use `> file 2>&1` instead.",
				Line:    redirectArg.TokenLiteralNode().Line,
				Column:  redirectArg.TokenLiteralNode().Column,
			},
		}
	}

	return nil
}
