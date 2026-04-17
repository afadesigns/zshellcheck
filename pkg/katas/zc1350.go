package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1350",
		Title:    "Use `${str:pos:len}` instead of `expr substr` for substring extraction",
		Severity: SeverityStyle,
		Description: "Zsh parameter expansion `${str:pos:len}` extracts a substring starting at " +
			"`pos` of length `len`. No external `expr` call, and the semantics are consistent " +
			"with `${str:pos}` (to end) and negative positions.",
		Check: checkZC1350,
	})
}

func checkZC1350(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "expr" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "substr" {
			return []Violation{{
				KataID: "ZC1350",
				Message: "Use `${str:pos:len}` instead of `expr substr` for substring extraction. " +
					"Parameter expansion avoids spawning an external process.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
