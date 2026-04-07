package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1269",
		Title:    "Use `pgrep` instead of `ps aux | grep` for process search",
		Severity: SeverityStyle,
		Description: "`ps aux | grep` matches itself in the process list requiring workarounds. " +
			"Use `pgrep` which is designed for process searching without self-matching.",
		Check: checkZC1269,
	})
}

func checkZC1269(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ps" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "aux" || val == "-ef" || val == "-e" {
			return []Violation{{
				KataID: "ZC1269",
				Message: "Use `pgrep` instead of `ps " + val + " | grep`. `pgrep` is purpose-built " +
					"for process searching and doesn't match itself.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
