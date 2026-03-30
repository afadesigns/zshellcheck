package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1171",
		Title:    "Use `print` instead of `echo -e` for escape sequences",
		Severity: SeverityStyle,
		Description: "`echo -e` behavior varies across shells and platforms. " +
			"In Zsh, `print` natively interprets escape sequences and is more reliable.",
		Check: checkZC1171,
	})
}

func checkZC1171(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "echo" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}

	first := cmd.Arguments[0].String()
	if first == "-e" {
		return []Violation{{
			KataID: "ZC1171",
			Message: "Use `print` instead of `echo -e`. Zsh `print` natively interprets " +
				"escape sequences and is more portable than `echo -e`.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
