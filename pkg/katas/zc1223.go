// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1223",
		Title:    "Avoid `ip addr show` piped to `grep` — use `ip -br addr`",
		Severity: SeverityStyle,
		Description: "`ip addr show | grep` parses verbose output. " +
			"`ip -br addr` provides machine-readable brief output without needing grep.",
		Check: checkZC1223,
	})
}

func checkZC1223(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ip" {
		return nil
	}

	hasAddr := false
	hasBrief := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "addr" || val == "address" {
			hasAddr = true
		}
		if val == "-br" || val == "-brief" {
			hasBrief = true
		}
	}

	if hasAddr && !hasBrief {
		return []Violation{{
			KataID: "ZC1223",
			Message: "Use `ip -br addr` for machine-readable output instead of " +
				"parsing `ip addr show` with grep or awk.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
