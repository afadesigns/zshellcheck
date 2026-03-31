package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1210",
		Title:    "Use `journalctl --no-pager` in scripts",
		Severity: SeverityStyle,
		Description: "`journalctl` invokes a pager by default which hangs in non-interactive scripts. " +
			"Use `--no-pager` for reliable script output.",
		Check: checkZC1210,
	})
}

func checkZC1210(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "journalctl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "--no-pager" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1210",
		Message: "Use `journalctl --no-pager` in scripts. Without it, " +
			"journalctl invokes a pager that hangs in non-interactive execution.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
