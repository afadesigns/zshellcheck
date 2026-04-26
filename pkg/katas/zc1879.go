// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1879",
		Title:    "Warn on `unsetopt BAD_PATTERN` — malformed glob patterns silently pass through as literals",
		Severity: SeverityWarning,
		Description: "`BAD_PATTERN` is on in Zsh by default: a syntactically broken glob (unbalanced " +
			"`[`, stray `^` outside extended-glob context, runaway `(alt|…`) produces a " +
			"`zsh: bad pattern` error so the script knows the filename filter is wrong. " +
			"Turning the option off reverts to POSIX behaviour — the pattern is handed to " +
			"the command verbatim, and `rm [abc` silently tries to remove a file literally " +
			"called `[abc`. Malformed patterns routed to `find -name` or passed to `case` " +
			"blocks likewise stop firing. Keep the option on at script level; if one " +
			"particular line really needs POSIX pass-through, quote the pattern or scope " +
			"with `setopt LOCAL_OPTIONS; unsetopt BAD_PATTERN` inside a function.",
		Check: checkZC1879,
	})
}

func checkZC1879(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			if zc1879IsBadPattern(arg.String()) {
				return zc1879Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOBADPATTERN" {
				return zc1879Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1879IsBadPattern(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "BADPATTERN"
}

func zc1879Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1879",
		Message: "`" + where + "` silences `bad pattern` errors — `rm [abc` tries " +
			"to remove a literal `[abc`, broken `case` arms stop firing. Keep " +
			"the option on; quote one-off patterns or scope with `LOCAL_OPTIONS`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
