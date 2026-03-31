package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1202",
		Title:    "Avoid `ifconfig` — use `ip` for network configuration",
		Severity: SeverityInfo,
		Description: "`ifconfig` is deprecated on modern Linux. " +
			"Use `ip addr`, `ip link`, or `ip route` from iproute2 for network operations.",
		Check: checkZC1202,
	})
}

func checkZC1202(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ifconfig" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1202",
		Message: "Avoid `ifconfig` — it is deprecated on modern Linux. " +
			"Use `ip addr`, `ip link`, or `ip route` from iproute2.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}
