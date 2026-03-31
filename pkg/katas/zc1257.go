package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1257",
		Title:    "Use `docker stop -t` to set graceful shutdown timeout",
		Severity: SeverityStyle,
		Description: "`docker stop` defaults to 10s before SIGKILL. In CI scripts, " +
			"set an explicit timeout with `-t` to control shutdown behavior.",
		Check: checkZC1257,
	})
}

func checkZC1257(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "stop" {
		return nil
	}

	hasTimeout := false
	for _, arg := range cmd.Arguments[1:] {
		if arg.String() == "-t" {
			hasTimeout = true
		}
	}

	if !hasTimeout {
		return []Violation{{
			KataID: "ZC1257",
			Message: "Use `docker stop -t N` to set an explicit shutdown timeout. " +
				"The default 10s may be too long or too short for your use case.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
