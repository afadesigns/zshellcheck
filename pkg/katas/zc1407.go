// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1407",
		Title:    "Avoid `/dev/tcp/...` — use Zsh `zsh/net/tcp` module",
		Severity: SeverityError,
		Description: "`/dev/tcp/host/port` is a Bash-specific virtual-file interface for TCP " +
			"connections; Zsh does not implement it. For TCP in Zsh, load `zmodload zsh/net/tcp` " +
			"and use `ztcp host port` which exposes the connection as a regular file descriptor.",
		Check: checkZC1407,
	})
}

func checkZC1407(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check all args for /dev/tcp or /dev/udp paths
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "/dev/tcp/") || strings.Contains(v, "/dev/udp/") {
			return []Violation{{
				KataID: "ZC1407",
				Message: "`/dev/tcp/...` and `/dev/udp/...` are Bash-only virtual files. In Zsh " +
					"load `zsh/net/tcp` and use `ztcp host port` / `ztcp -c $fd` for TCP.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}
