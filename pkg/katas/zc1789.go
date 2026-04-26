// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1789",
		Title:    "Warn on `setopt CORRECT` / `CORRECT_ALL` — Zsh spellcheck silently rewrites script tokens",
		Severity: SeverityWarning,
		Description: "`setopt CORRECT` prompts to rewrite command names that look mistyped; " +
			"`CORRECT_ALL` extends the check to every argument on the line. In an interactive " +
			"shell this is a friendly nudge. In a script it becomes a footgun: a filename " +
			"that is *close enough* to an existing file gets silently replaced with that " +
			"other file, and the \"nlh?\" prompt reads from stdin — which may be the input " +
			"the script was supposed to process. Keep `CORRECT` / `CORRECT_ALL` in " +
			"`~/.zshrc` only and never toggle them inside a function a script calls.",
		Check: checkZC1789,
	})
}

func checkZC1789(node ast.Node) []Violation {
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
			if name, hit := zc1789Matches(arg.String()); hit {
				return zc1789Hit(cmd, "setopt "+arg.String(), name)
			}
		}
	case "set":
		for i, arg := range cmd.Arguments {
			v := arg.String()
			if (v == "-o" || v == "--option") && i+1 < len(cmd.Arguments) {
				next := cmd.Arguments[i+1].String()
				if name, hit := zc1789Matches(next); hit {
					return zc1789Hit(cmd, "set -o "+next, name)
				}
			}
		}
	}
	return nil
}

func zc1789Matches(v string) (string, bool) {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	switch norm {
	case "CORRECT":
		return "CORRECT", true
	case "CORRECTALL":
		return "CORRECT_ALL", true
	}
	return "", false
}

func zc1789Hit(cmd *ast.SimpleCommand, where, canonical string) []Violation {
	return []Violation{{
		KataID: "ZC1789",
		Message: "`" + where + "` enables `" + canonical + "` — Zsh spellcheck " +
			"silently rewrites tokens that look mistyped. In a script that corrupts " +
			"file paths and steals stdin for the correction prompt. Keep in `~/.zshrc`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
