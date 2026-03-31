package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1192",
		Title:    "Avoid `sleep 0` — it is a no-op external process",
		Severity: SeverityInfo,
		Description: "`sleep 0` spawns an external process that does nothing. " +
			"Remove it or use `:` if an explicit no-op is needed.",
		Check: checkZC1192,
	})
}

func checkZC1192(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sleep" {
		return nil
	}

	if len(cmd.Arguments) == 1 && cmd.Arguments[0].String() == "0" {
		return []Violation{{
			KataID: "ZC1192",
			Message: "Remove `sleep 0` — it spawns a process that does nothing. " +
				"Use `:` if an explicit no-op is needed.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityInfo,
		}}
	}

	return nil
}
