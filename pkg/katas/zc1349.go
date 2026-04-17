package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1349",
		Title:    "Use `${#var}` instead of `expr length \"$var\"` for string length",
		Severity: SeverityStyle,
		Description: "Zsh (and POSIX) `${#var}` returns string length without spawning `expr`. " +
			"Use it wherever you would reach for `expr length` or `expr STRING : '.*'`.",
		Check: checkZC1349,
	})
}

func checkZC1349(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "expr" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "length" {
			return []Violation{{
				KataID: "ZC1349",
				Message: "Use `${#var}` instead of `expr length \"$var\"` for string length. " +
					"Parameter expansion avoids spawning an external process.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
