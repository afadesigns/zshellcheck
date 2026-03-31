package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1215",
		Title:    "Source `/etc/os-release` instead of parsing with `cat`/`grep`",
		Severity: SeverityStyle,
		Description: "`/etc/os-release` is designed to be sourced directly. " +
			"Use `. /etc/os-release` to get variables like `$ID`, `$VERSION_ID` without parsing.",
		Check: checkZC1215,
	})
}

func checkZC1215(node ast.Node) []Violation {
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
		if val == "/etc/os-release" || val == "/etc/lsb-release" {
			return []Violation{{
				KataID: "ZC1215",
				Message: "Source `/etc/os-release` directly with `. /etc/os-release` instead of " +
					"parsing with `cat`. It exports variables like `$ID` and `$VERSION_ID`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
