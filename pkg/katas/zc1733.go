// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1733",
		Title:    "Error on `docker plugin install --grant-all-permissions` — accepts every requested cap",
		Severity: SeverityError,
		Description: "Docker plugins run as root with whatever privileges they ask for at install " +
			"time — host networking, `/dev/*` mounts, arbitrary capability grants. The " +
			"interactive prompt enumerates each request so the operator can refuse anything " +
			"unexpected. `--grant-all-permissions` skips the prompt and accepts the whole " +
			"list, so a compromised plugin author or a typo-squatted name owns the host " +
			"on first install. Install plugins by name, walk the prompt manually, then pin " +
			"the tag (`@sha256:...`) once vetted.",
		Check: checkZC1733,
	})
}

func checkZC1733(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}
	if len(cmd.Arguments) < 3 {
		return nil
	}
	if cmd.Arguments[0].String() != "plugin" || cmd.Arguments[1].String() != "install" {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		if arg.String() == "--grant-all-permissions" {
			return []Violation{{
				KataID: "ZC1733",
				Message: "`docker plugin install --grant-all-permissions` accepts every " +
					"capability the plugin requests — root-equivalent on the host. Walk " +
					"the interactive prompt manually and pin the digest once vetted.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
