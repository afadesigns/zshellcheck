// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1599",
		Title:    "Warn on `ldconfig -f PATH` outside `/etc/` — attacker-writable loader cache",
		Severity: SeverityWarning,
		Description: "`ldconfig -f PATH` rebuilds `/etc/ld.so.cache` using PATH instead of the " +
			"system `/etc/ld.so.conf`. If PATH sits in `/tmp`, `/var/tmp`, `$HOME`, or any " +
			"directory an attacker can create, they can inject an `include` line that points " +
			"at their directory of malicious shared objects. After the cache rebuild, every " +
			"subsequent executable on the host loads their library first. Keep the config " +
			"under `/etc/ld.so.conf.d/` with root ownership and run `ldconfig` with no `-f`.",
		Check: checkZC1599,
	})
}

func checkZC1599(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ldconfig" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		if arg.String() != "-f" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		next := cmd.Arguments[i+1].String()
		if strings.HasPrefix(next, "/etc/") || strings.HasPrefix(next, "$") {
			continue
		}
		return []Violation{{
			KataID: "ZC1599",
			Message: "`ldconfig -f " + next + "` uses a config outside `/etc/`. If the " +
				"file is attacker-writable, every binary on the host loads the attacker's " +
				"library. Keep config under `/etc/ld.so.conf.d/`.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}
