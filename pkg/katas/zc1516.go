// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1516",
		Title:    "Error on `umask 000` / `umask 0` — new files / directories world-writable",
		Severity: SeverityError,
		Description: "`umask 000` means every file created after this line inherits mode 0666 " +
			"and every directory inherits 0777 — world-readable, world-writable, no " +
			"authorization layer. On a multi-user host (build runner, shared workstation) this " +
			"leaks secrets through the filesystem and invites tampering. Pick a sensible " +
			"umask (`022` for public software, `077` for secrets handling).",
		Check: checkZC1516,
	})
}

func checkZC1516(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "umask" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}
	v := cmd.Arguments[0].String()
	if v == "0" || v == "00" || v == "000" || v == "0000" {
		return []Violation{{
			KataID: "ZC1516",
			Message: "`umask " + v + "` leaves new files world-readable and world-writable. " +
				"Use `022` for public software, `077` for secrets handling.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}
	return nil
}
