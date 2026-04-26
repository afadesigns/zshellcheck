// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1832",
		Title:    "Warn on Zsh `limit coredumpsize unlimited` — setuid memory landing in core files",
		Severity: SeverityWarning,
		Description: "Zsh's `limit` builtin is the csh-style sibling of `ulimit`; `limit " +
			"coredumpsize unlimited` is the Zsh equivalent of `ulimit -c unlimited` and has " +
			"the same consequence: a crashing setuid or key-holding process leaves its " +
			"address space on disk as a world-readable core file. Leave the coredump " +
			"ceiling at the distro default (usually 0 for non-debug sessions), or use " +
			"`systemd-coredump` with restricted permissions when you need post-mortem data. " +
			"`ulimit -c unlimited` is covered by ZC1495; this kata catches the Zsh-specific " +
			"`limit`/`unlimit coredumpsize` spelling.",
		Check: checkZC1832,
	})
}

func checkZC1832(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "limit":
		// `limit coredumpsize unlimited` (with optional -h for hard limit).
		argIdx := 0
		if argIdx < len(cmd.Arguments) && cmd.Arguments[argIdx].String() == "-h" {
			argIdx++
		}
		if argIdx+1 >= len(cmd.Arguments) {
			return nil
		}
		resource := strings.ToLower(cmd.Arguments[argIdx].String())
		value := strings.ToLower(cmd.Arguments[argIdx+1].String())
		if (resource == "coredumpsize" || resource == "coredump") && value == "unlimited" {
			return zc1832Hit(cmd, "limit coredumpsize unlimited")
		}
	case "unlimit":
		for _, arg := range cmd.Arguments {
			v := strings.ToLower(arg.String())
			if v == "coredumpsize" || v == "coredump" {
				return zc1832Hit(cmd, "unlimit coredumpsize")
			}
		}
	}
	return nil
}

func zc1832Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1832",
		Message: "`" + where + "` enables unbounded core dumps (Zsh-specific `limit` " +
			"spelling of `ulimit -c unlimited`). A setuid crash drops its memory to " +
			"disk as a world-readable file — leave the ceiling at the distro default.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
