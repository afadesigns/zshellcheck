package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1287",
		Title:    "Use `cat -v` alternative: Zsh `${(V)var}` for visible control characters",
		Severity: SeverityStyle,
		Description: "Zsh provides the `(V)` parameter expansion flag to make control characters " +
			"visible in a variable. This avoids piping through `cat -v` for simple " +
			"visibility of non-printable characters.",
		Check: checkZC1287,
	})
}

func checkZC1287(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cat" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-v" || val == "-A" {
			return []Violation{{
				KataID:  "ZC1287",
				Message: "Use Zsh `${(V)var}` to make control characters visible instead of `cat -v`.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityStyle,
			}}
		}
	}

	return nil
}
