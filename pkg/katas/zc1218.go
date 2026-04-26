// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1218",
		Title:    "Avoid `useradd` without `--shell /sbin/nologin` for service accounts",
		Severity: SeverityWarning,
		Description: "Service accounts created with `useradd` should use `--shell /sbin/nologin` " +
			"and `--system` to prevent interactive login and use system UID ranges.",
		Check: checkZC1218,
	})
}

func checkZC1218(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "useradd" {
		return nil
	}

	hasSystem := false
	hasNologin := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "--system" || val == "-r" {
			hasSystem = true
		}
		if val == "/sbin/nologin" || val == "/usr/sbin/nologin" || val == "/bin/false" {
			hasNologin = true
		}
	}

	if hasSystem && !hasNologin {
		return []Violation{{
			KataID: "ZC1218",
			Message: "Add `--shell /sbin/nologin` when creating system accounts with `useradd`. " +
				"This prevents interactive login for service accounts.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
