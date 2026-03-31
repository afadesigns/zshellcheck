package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1205",
		Title:    "Avoid `arp` — use `ip neigh` for neighbor tables",
		Severity: SeverityInfo,
		Description: "`arp` is deprecated on modern Linux in favor of `ip neigh` from iproute2. " +
			"`ip neigh` provides consistent syntax with other `ip` subcommands.",
		Check: checkZC1205,
	})
}

func checkZC1205(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "arp" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1205",
		Message: "Avoid `arp` — it is deprecated on modern Linux. " +
			"Use `ip neigh` from iproute2 for neighbor table management.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}
