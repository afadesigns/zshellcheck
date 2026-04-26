// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1604",
		Title:    "Warn on `source <glob>` / `. <glob>` — loads every match; one bad file = code exec",
		Severity: SeverityWarning,
		Description: "`source /etc/profile.d/*.sh` and similar glob-sourcing patterns load every " +
			"file that matches, in the order Zsh enumerates them. One attacker-writable file " +
			"anywhere in the glob yields arbitrary code execution as whoever is running the " +
			"script, with that caller's privileges. Prefer explicit filenames so review can " +
			"enumerate exactly what gets loaded. If a directory of drop-ins is required, audit " +
			"ownership and perms at install time and keep the directory root-owned.",
		Check: checkZC1604,
	})
}

func checkZC1604(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "source" && ident.Value != "." {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}

	target := cmd.Arguments[0].String()
	if !strings.ContainsAny(target, "*?[") {
		return nil
	}

	return []Violation{{
		KataID: "ZC1604",
		Message: "`" + ident.Value + " " + target + "` loads every matched file. One " +
			"attacker-writable match is arbitrary code execution. Use explicit filenames.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
