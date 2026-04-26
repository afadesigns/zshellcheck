// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1894",
		Title:    "Error on `conntrack -F` / `--flush` — every tracked connection (including SSH) is reset",
		Severity: SeverityError,
		Description: "`conntrack -F` (alias `--flush`) wipes the netfilter connection-tracking " +
			"table. Every established TCP flow that depended on conntrack (every " +
			"stateful-NAT connection, every `-m conntrack --ctstate RELATED,ESTABLISHED` " +
			"allowance, every MASQUERADE session) loses its entry and the next packet is " +
			"matched from scratch; most firewall rulesets drop it as \"new\" and the " +
			"session dies. Over SSH, that means the shell running the very command drops. " +
			"Stage the flush behind `at now + 5 minutes` so the session can re-enter the " +
			"table via a preceding rule, or narrow the scope with `conntrack -D -s " +
			"<client-IP>` for a specific hung flow.",
		Check: checkZC1894,
	})
}

var zc1894FlushFlags = map[string]bool{
	"-F":      true,
	"--flush": true,
}

func checkZC1894(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "conntrack" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if zc1894FlushFlags[arg.String()] {
			line, col := FlagArgPosition(cmd, zc1894FlushFlags)
			return []Violation{{
				KataID: "ZC1894",
				Message: "`conntrack -F` wipes every tracked flow — stateful " +
					"`ctstate ESTABLISHED` allowances drop, running SSH sessions " +
					"lose their entry. Gate with `at now + N min` or narrow to " +
					"one flow with `conntrack -D -s <ip>`.",
				Line:   line,
				Column: col,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
