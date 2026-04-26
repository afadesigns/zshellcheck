// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1838",
		Title:    "Warn on `setopt GLOB_DOTS` — bare `*` silently starts matching hidden files",
		Severity: SeverityWarning,
		Description: "`GLOB_DOTS` off is the Zsh default: patterns like `*`, `*.log`, and " +
			"recursive `**/*` skip filenames that begin with a dot (`.git/`, `.env`, " +
			"`.ssh/`). Setting `setopt GLOB_DOTS` script-wide reverses that quietly — every " +
			"subsequent glob now also matches hidden entries, which turns routine " +
			"maintenance lines (`rm *`, `cp -r * /backup`, `chmod 644 *`) into " +
			"repository-wiping, secret-copying, permission-flipping bugs. Leave the option " +
			"alone at the script level and request dot-inclusion per-glob with the " +
			"Zsh-native `*(D)` qualifier (or `.* *` when you explicitly want both), so the " +
			"effect is scoped to the exact line that needs it.",
		Check: checkZC1838,
	})
}

func checkZC1838(node ast.Node) []Violation {
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
			v := arg.String()
			if zc1838IsGlobDots(v) {
				return zc1838Hit(cmd, "setopt "+v)
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOGLOBDOTS" {
				return zc1838Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1838IsGlobDots(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "GLOBDOTS"
}

func zc1838Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1838",
		Message: "`" + where + "` makes every bare `*` also match hidden files — " +
			"`rm *` quietly destroys `.git/`, `cp -r *` copies `.env`. Keep the " +
			"option alone; request dotfiles per-glob with `*(D)`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
