// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1455",
		Title:    "Avoid `docker run --net=host` / `--network=host` — disables network isolation",
		Severity: SeverityWarning,
		Description: "Host networking gives the container direct access to the host's network " +
			"stack, including localhost services. A vulnerable container can reach services " +
			"meant to be local-only. Use `-p hostport:containerport` for specific publishes and " +
			"dedicated networks for inter-container traffic.",
		Check: checkZC1455,
	})
}

func checkZC1455(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "docker" && ident.Value != "podman" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--net=host" || v == "--network=host" || v == "-net=host" {
			return []Violation{{
				KataID: "ZC1455",
				Message: "`--net=host` / `--network=host` lets the container reach host-local " +
					"services. Use `-p` for explicit publishes or dedicated container networks.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
