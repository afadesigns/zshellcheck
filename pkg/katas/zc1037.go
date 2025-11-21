package katas

import (

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/token"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1037",
		Title:       "Use 'print -r --' for variable expansion",
		Description: "Using 'echo' to print strings containing variables can lead to unexpected behavior " +
			"if the variable contains special characters or flags. A safer, more reliable alternative " +
			"is 'print -r --'.",
		Check:       checkZC1037,
	})
}

func checkZC1037(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if cmd.Name.TokenLiteral() != "echo" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if ident, ok := arg.(*ast.Identifier); ok && ident.Token.Type == token.VARIABLE {
			return []Violation{
				{
					KataID:  "ZC1037",
					Message: "Use 'print -r --' instead of 'echo' to reliably print variable expansions.",
					Line:    cmd.Token.Line,
					Column:  cmd.Token.Column,
				},
			}
		}
	}

	return nil
}