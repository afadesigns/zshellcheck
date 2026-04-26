// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1885",
		Title:    "Warn on `setopt CSH_NULL_GLOB` — unmatched globs drop instead of erroring when any sibling matches",
		Severity: SeverityWarning,
		Description: "`CSH_NULL_GLOB` (off by default) mimics csh's rule: in a list like " +
			"`rm *.log *.bak *.tmp`, if at least one pattern produces matches the " +
			"remaining unmatched patterns are silently discarded, and only if every " +
			"pattern produces nothing does the shell raise `no match`. That is a " +
			"partial-failure concealer — a genuine typo `rm *.lg *.bak` can still " +
			"delete the `.bak` files while hiding the `.lg` mismatch, and maintenance " +
			"loops that relied on `NOMATCH` to stop on typos pass right through. Keep " +
			"the option off at script level; use `*(N)` per-glob when you want " +
			"null-glob behaviour.",
		Check: checkZC1885,
	})
}

func checkZC1885(node ast.Node) []Violation {
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
			if zc1885IsCshNullGlob(arg.String()) {
				return zc1885Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOCSHNULLGLOB" {
				return zc1885Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1885IsCshNullGlob(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "CSHNULLGLOB"
}

func zc1885Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1885",
		Message: "`" + where + "` silently discards unmatched globs in a list when " +
			"any sibling matches — `rm *.lg *.bak` deletes the `.bak` files " +
			"and hides the typo. Keep the option off; use `*(N)` per-glob.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
