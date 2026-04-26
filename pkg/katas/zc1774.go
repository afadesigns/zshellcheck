// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1774",
		Title:    "Warn on `setopt GLOB_SUBST` — `$var` starts glob-expanding, user data becomes a pattern",
		Severity: SeverityWarning,
		Description: "With `GLOB_SUBST` enabled, the result of any parameter expansion is " +
			"rescanned for filename-generation metacharacters (`*`, `?`, `[`, `^`, `~`, " +
			"brace ranges, qualifiers). Zsh's default — `NO_GLOB_SUBST` — keeps `$var` literal " +
			"and matches the behavior most script authors expect after moving from Bash or " +
			"POSIX sh. Turning `GLOB_SUBST` on globally means any unquoted `$var` that " +
			"contains a metacharacter (environment, argv, file contents, user prompt) is " +
			"expanded against the filesystem — an injection vector, and a subtle source of " +
			"`no matches found` failures on empty variables. Keep `setopt GLOB_SUBST` inside a " +
			"narrow subshell or function body, or use explicit `~` / `(e)` / `(P)` flags where " +
			"you actually want the rescan.",
		Check: checkZC1774,
	})
}

func checkZC1774(node ast.Node) []Violation {
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
			v := strings.ToUpper(strings.ReplaceAll(arg.String(), "_", ""))
			if v == "GLOBSUBST" {
				return zc1774Hit(cmd, "setopt "+arg.String())
			}
		}
	case "set":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "-G" {
				return zc1774Hit(cmd, "set -G")
			}
			if v == "-o" || v == "--option" {
				continue
			}
			if strings.EqualFold(strings.ReplaceAll(v, "_", ""), "globsubst") {
				return zc1774Hit(cmd, "set -o "+v)
			}
		}
	}
	return nil
}

func zc1774Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1774",
		Message: "`" + where + "` enables `GLOB_SUBST` — every unquoted `$var` expansion " +
			"is rescanned as a glob pattern. User-controlled data becomes a filesystem " +
			"query. Scope this in a subshell / function, or use explicit expansion flags.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
