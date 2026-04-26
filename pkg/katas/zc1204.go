// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1204",
		Title:    "Avoid `route` — use `ip route` for routing",
		Severity: SeverityInfo,
		Description: "`route` is deprecated on modern Linux in favor of `ip route` from iproute2. " +
			"`ip route` provides consistent syntax with other `ip` subcommands.",
		Check: checkZC1204,
	})
}

func checkZC1204(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "route" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1204",
		Message: "Avoid `route` — it is deprecated on modern Linux. " +
			"Use `ip route` from iproute2 for consistent routing management.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}
