package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1121",
		Title: "Use `$HOST` instead of `hostname`",
		Description: "Zsh provides `$HOST` as a built-in variable containing the hostname. " +
			"Avoid spawning `hostname` as an external process.",
		Severity: SeverityStyle,
		Check:    checkZC1121,
	})
}

func checkZC1121(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "hostname" {
		return nil
	}

	// hostname with flags like -f, -I, -d does more than $HOST
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1121",
		Message: "Use `$HOST` instead of `hostname`. " +
			"Zsh maintains `$HOST` as a built-in variable, avoiding an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
