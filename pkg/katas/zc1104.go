package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1104",
		Title: "Suggest `path` array instead of `export PATH` string manipulation",
		Description: "Zsh automatically maps the `$PATH` environment variable to the `$path` array. " +
			"Modifying `$path` is cleaner and less error-prone than manipulating the colon-separated `$PATH` string.",
		Check: checkZC1104,
	})
}

func checkZC1104(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check for `export PATH=...`
	if cmd.Name != nil && cmd.Name.String() == "export" {
		for _, arg := range cmd.Arguments {
			argStr := arg.String()
			if strings.HasPrefix(argStr, "PATH=") {
				return []Violation{{
					KataID:  "ZC1104",
					Message: "Use the `path` array instead of manually manipulating the `$PATH` string.",
					Line:    cmd.Token.Line,
					Column:  cmd.Token.Column,
				}}
			}
		}
	}

	return nil
}
