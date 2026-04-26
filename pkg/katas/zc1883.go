// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1883",
		Title:    "Warn on `setopt PATH_SCRIPT` — `. ./script.sh` silently falls back to `$PATH` lookup",
		Severity: SeverityWarning,
		Description: "`PATH_SCRIPT` (off by default) lets the `.` builtin and `source` fall back to " +
			"a `$PATH` walk when the literal path resolves to no file. With it on, " +
			"`. helper.sh` looks for `helper.sh` in every `$path` entry — including " +
			"user-owned directories like `~/bin` or `./` — and silently sources whichever " +
			"matches first. An attacker who can drop `helper.sh` into any `$PATH` " +
			"component runs their code inside the current shell's process, with every " +
			"parent env var and exported secret available. Keep the option off; always " +
			"source scripts with an explicit path (`./helper.sh`, `/opt/…/helper.sh`) so " +
			"the source cannot be redirected.",
		Check: checkZC1883,
	})
}

func checkZC1883(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			if zc1883IsPathScript(arg.String()) {
				return zc1883Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOPATHSCRIPT" {
				return zc1883Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1883IsPathScript(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "PATHSCRIPT"
}

func zc1883Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1883",
		Message: "`" + where + "` lets `.`/`source` fall back to `$PATH` when a " +
			"literal path misses — a dropper in `~/bin` or `./` runs inside the " +
			"current shell with every exported secret. Keep the option off; " +
			"always use explicit paths.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
