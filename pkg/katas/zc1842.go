// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1842",
		Title:    "Warn on `setopt CDABLE_VARS` — `cd NAME` silently falls back to `cd $NAME`",
		Severity: SeverityWarning,
		Description: "With `CDABLE_VARS` on, any `cd NAME` whose `NAME` does not exist as a " +
			"directory is retried as `cd ${NAME}` — if a parameter of the same name is set, " +
			"the working directory silently jumps to wherever the variable points. A typo " +
			"like `cd cinfig` (intent: `config`) suddenly lands inside `${cinfig}` when one " +
			"exists, and every later relative path in the script is computed from the wrong " +
			"root. Keep this option inside `~/.zshrc` where it is an interactive shortcut; " +
			"in scripts, always `cd \"$dir\"` explicitly and pair with `|| exit` so a missed " +
			"directory fails loudly instead of rewriting `$PWD`.",
		Check: checkZC1842,
	})
}

func checkZC1842(node ast.Node) []Violation {
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
			if zc1842IsCdableVars(arg.String()) {
				return zc1842Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOCDABLEVARS" {
				return zc1842Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1842IsCdableVars(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "CDABLEVARS"
}

func zc1842Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1842",
		Message: "`" + where + "` turns a failed `cd NAME` into `cd $NAME` — a typo " +
			"silently lands in whatever directory the matching variable points to. " +
			"Keep this in `~/.zshrc`; in scripts use `cd \"$dir\" || exit`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
