// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1225",
		Title:    "Avoid parsing `uptime` — read `/proc/uptime` directly",
		Severity: SeverityStyle,
		Description: "`uptime` output is human-readable and varies by locale. " +
			"Read `/proc/uptime` for machine-parseable uptime in seconds.",
		Check: checkZC1225,
	})
}

func checkZC1225(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "uptime" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1225",
		Message: "Avoid parsing `uptime` — its output varies by locale. " +
			"Read `/proc/uptime` for machine-parseable seconds since boot.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
