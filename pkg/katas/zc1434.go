package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1434",
		Title:    "Warn on `swapoff -a` — disables all swap, can OOM-kill",
		Severity: SeverityWarning,
		Description: "`swapoff -a` disables every active swap. On a memory-constrained host " +
			"this pushes data back into RAM, potentially triggering OOM-killer. Prefer " +
			"disabling specific devices/files (`swapoff /swapfile`) and verify memory headroom " +
			"with `free -m` first.",
		Check: checkZC1434,
	})
}

func checkZC1434(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "swapoff" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-a" || arg.String() == "--all" {
			return []Violation{{
				KataID: "ZC1434",
				Message: "`swapoff -a` disables ALL swap areas — risks OOM on memory-constrained " +
					"hosts. Disable specific swaps (`swapoff /swapfile`) after checking `free -m`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
