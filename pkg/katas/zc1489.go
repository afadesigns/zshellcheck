// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1489",
		Title:    "Error on `nc -e` / `ncat -e` — classic reverse-shell invocation",
		Severity: SeverityError,
		Description: "`nc -e <shell>` and `ncat -e <shell>` pipe a shell to a network socket. " +
			"This is the canonical reverse-shell payload. Most distro builds of `nc` have " +
			"`-e` disabled for precisely this reason, so seeing it in a script is either an " +
			"attacker backdoor or a deployment time bomb waiting on a different packaging " +
			"of netcat. If you need a bidirectional pipe, use `socat TCP:... EXEC:...,pty` " +
			"with an explicit authorization check and document the use.",
		Check: checkZC1489,
	})
}

func checkZC1489(node ast.Node) []Violation {
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
				KataID: "ZC1489",
				Message: "`" + ident.Value + " " + v + "` is the classic reverse-shell flag. " +
					"Use socat with explicit PTY + authorization instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
