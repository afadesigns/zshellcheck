package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1102",
		Title: "Redirecting output of `sudo` doesn't work as expected",
		Description: "Redirections are performed by the current shell before `sudo` is started. " +
			"So `sudo echo > /root/file` will try to open `/root/file` as the current user, failing. " +
			"Use `echo ... | sudo tee file` or `sudo sh -c 'echo ... > file'`.",
		Check: checkZC1102,
	})
}

func checkZC1102(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if the command name is 'sudo'
	if cmd.Name != nil && cmd.Name.String() == "sudo" {
		// Scan arguments for output redirection operators
		for _, arg := range cmd.Arguments {
			argStr := arg.String()
			if argStr == ">" || argStr == ">>" {
				return []Violation{{
					KataID:  "ZC1102",
					Message: "Redirection happens before `sudo`. This will likely fail permission checks. Use `| sudo tee`.",
					Line:    cmd.TokenLiteralNode().Line,
					Column:  cmd.TokenLiteralNode().Column,
				}}
			}
		}
	}

	return nil
}
