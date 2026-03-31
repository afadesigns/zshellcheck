package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1219",
		Title:    "Use `curl -fsSL` instead of `wget -O -` for piped downloads",
		Severity: SeverityStyle,
		Description: "`wget -O -` outputs to stdout but lacks `curl`'s error handling. " +
			"`curl -fsSL` fails on HTTP errors, is silent, follows redirects, and is more portable.",
		Check: checkZC1219,
	})
}

func checkZC1219(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "wget" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-O-" || val == "-qO-" {
			return []Violation{{
				KataID: "ZC1219",
				Message: "Use `curl -fsSL` instead of `wget -O -` for piped downloads. " +
					"`curl` fails on HTTP errors and is available on more platforms.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
