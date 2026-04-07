package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1272",
		Title:    "Use `install -m` instead of separate `cp` and `chmod`",
		Severity: SeverityStyle,
		Description: "`install` atomically copies a file and sets permissions in one step. " +
			"Using separate `cp` and `chmod` creates a window where the file has wrong permissions.",
		Check: checkZC1272,
	})
}

func checkZC1272(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cp" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "/usr/local/bin" || val == "/usr/bin" || val == "/opt/bin" ||
			val == "/usr/local/sbin" || val == "/usr/sbin" {
			return []Violation{{
				KataID:  "ZC1272",
				Message: "Use `install -m 0755` instead of `cp` to system directories. `install` sets permissions atomically.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityStyle,
			}}
		}
	}

	return nil
}
