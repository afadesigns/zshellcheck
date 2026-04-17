package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1422",
		Title:    "Avoid `sudo -S` — reads password from stdin, exposes plaintext",
		Severity: SeverityError,
		Description: "`sudo -S` reads the password from stdin, enabling `echo $PW | sudo -S cmd` " +
			"patterns that place the plaintext password in the process tree and shell history. " +
			"Prefer `sudo -A` with a graphical askpass, `NOPASSWD:` in sudoers for specific " +
			"commands, or `pkexec` for policy-based privilege elevation.",
		Check: checkZC1422,
	})
}

func checkZC1422(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sudo" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-S" {
			return []Violation{{
				KataID: "ZC1422",
				Message: "`sudo -S` enables password-via-stdin. Avoid piping plaintext " +
					"credentials. Use `sudo -A` (askpass), `NOPASSWD:` in sudoers, or `pkexec`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}
