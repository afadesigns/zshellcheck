package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1238",
		Title:    "Avoid `docker exec -it` in scripts — drop `-it` for non-interactive",
		Severity: SeverityWarning,
		Description: "`docker exec -it` allocates a TTY and attaches stdin, which hangs " +
			"in non-interactive scripts. Use `docker exec` without `-it` for scripted commands.",
		Check: checkZC1238,
	})
}

func checkZC1238(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "exec" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "-it" || val == "-ti" {
			return []Violation{{
				KataID: "ZC1238",
				Message: "Avoid `docker exec -it` in scripts — TTY allocation hangs without a terminal. " +
					"Use `docker exec` without `-it` for non-interactive commands.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
