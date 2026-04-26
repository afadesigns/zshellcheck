// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1563",
		Title:    "Warn on `swapoff -a` — disables swap (memory pressure, potential OOM)",
		Severity: SeverityWarning,
		Description: "`swapoff -a` turns off every active swap device. Kubelet installers do " +
			"this because kubelet refuses to run with swap, but leaving it in a general-purpose " +
			"script means the next memory-hungry process on the host hits the OOM killer " +
			"instead of paging. If the goal is kubelet-friendly, also remove the swap entry " +
			"from `/etc/fstab` and document the trade-off; otherwise keep swap on.",
		Check: checkZC1563,
	})
}

func checkZC1563(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "swapoff" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-a" || arg.String() == "--all" {
			return []Violation{{
				KataID: "ZC1563",
				Message: "`swapoff -a` disables all swap devices — next memory-hungry process " +
					"hits OOM. Document the trade-off if kubelet requires it.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
