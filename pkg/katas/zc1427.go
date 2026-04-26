// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1427",
		Title:    "Dangerous: `nc -e` / `ncat -e` — spawns arbitrary command on network connect",
		Severity: SeverityError,
		Description: "`nc -e cmd` and `ncat --exec cmd` pipe the network socket to an arbitrary " +
			"command. Incoming connections get a shell or any command you specify — the " +
			"classic reverse-shell pattern. Many distros ship `nc` compiled without `-e` for " +
			"this reason. Remove `-e` from scripts except in audited, restricted contexts.",
		Check: checkZC1427,
	})
}

func checkZC1427(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "nc" && ident.Value != "ncat" && ident.Value != "netcat" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-e" || v == "-c" {
			return []Violation{{
				KataID: "ZC1427",
				Message: "`" + ident.Value + " " + v + "` spawns an arbitrary command for " +
					"each connection — reverse-shell territory. Remove from scripts unless " +
					"audited and restricted.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}
