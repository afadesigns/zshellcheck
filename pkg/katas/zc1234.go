package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1234",
		Title:    "Use `docker run --rm` to auto-remove containers",
		Severity: SeverityStyle,
		Description: "`docker run` without `--rm` leaves stopped containers behind. " +
			"Use `--rm` in scripts to automatically clean up after execution.",
		Check: checkZC1234,
	})
}

func checkZC1234(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "run" {
		return nil
	}

	hasRM := false
	hasDetach := false

	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "--rm" {
			hasRM = true
		}
		if val == "-d" {
			hasDetach = true
		}
	}

	if !hasRM && !hasDetach {
		return []Violation{{
			KataID: "ZC1234",
			Message: "Use `docker run --rm` to auto-remove containers after exit. " +
				"Without `--rm`, stopped containers accumulate on disk.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
