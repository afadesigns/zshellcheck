package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.RedirectionNode, Kata{
		ID:    "ZC1102",
		Title: "Redirecting output of `sudo` doesn't work as expected",
		Description: "Redirections are performed by the current shell before `sudo` is started. " +
			"So `sudo echo > /root/file` will try to open `/root/file` as the current user, failing. " +
			"Use `echo ... | sudo tee file` or `sudo sh -c 'echo ... > file'`.",
		Check: checkZC1102,
	})
}

func checkZC1102(node ast.Node) []Violation {
	redir, ok := node.(*ast.Redirection)
	if !ok {
		return nil
	}

	// Check if the left-hand side of the redirection is a SimpleCommand
	if cmd, ok := redir.Left.(*ast.SimpleCommand); ok {
		// Check if the command name is 'sudo'
		if cmd.Name != nil && cmd.Name.String() == "sudo" {
			// Check if it's an output redirection
			if redir.Operator == ">" || redir.Operator == ">>" {
				return []Violation{{
					KataID:  "ZC1102",
					Message: "Redirection happens before `sudo`. This will likely fail permission checks. Use `| sudo tee`.",
					Line:    redir.Token.Line,
					Column:  redir.Token.Column,
				}}
			}
		}
	}

	return nil
}
