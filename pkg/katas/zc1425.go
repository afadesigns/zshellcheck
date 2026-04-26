// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1425",
		Title:    "`shutdown` / `reboot` / `halt` / `poweroff` — confirm before scripting",
		Severity: SeverityWarning,
		Description: "Scripts that invoke `shutdown`, `reboot`, `halt`, `poweroff`, or " +
			"`systemctl poweroff` take down the system. Unattended invocation in automation is " +
			"often wrong (e.g. leftover test step). Prefer `systemctl isolate rescue.target` for " +
			"controlled scenarios, and require explicit confirmation for interactive scripts.",
		Check: checkZC1425,
	})
}

func checkZC1425(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "shutdown", "reboot", "halt", "poweroff":
		return []Violation{{
			KataID: "ZC1425",
			Message: "`" + ident.Value + "` takes down the system. In scripts, confirm the " +
				"caller really wants this (interactive prompt, feature flag, or CI guard).",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
