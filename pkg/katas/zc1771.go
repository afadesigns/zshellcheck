// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1771",
		Title:    "Warn on `alias -g` / `alias -s` — global and suffix aliases surprise script readers",
		Severity: SeverityWarning,
		Description: "`alias -g NAME=value` defines a global alias that expands anywhere on the " +
			"command line, not just in command position. `alias -s ext=cmd` (suffix alias) runs " +
			"`cmd file.ext` whenever a bare `file.ext` appears as a command. Both forms are " +
			"Zsh-idiomatic interactive conveniences; in scripts they produce surprising " +
			"substitutions that a reader cannot infer from local context — a bare word like " +
			"`G` or `foo.log` stops meaning what it looks like. Use a function or a regular " +
			"alias instead, and keep `alias -g` / `alias -s` in your `~/.zshrc` where the " +
			"definition is discoverable.",
		Check: checkZC1771,
	})
}

func checkZC1771(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "alias" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	first := cmd.Arguments[0].String()
	switch {
	case first == "-g":
		return zc1771Hit(cmd, "-g", "global")
	case first == "-s":
		return zc1771Hit(cmd, "-s", "suffix")
	case strings.HasPrefix(first, "-") && !strings.HasPrefix(first, "--"):
		if strings.ContainsRune(first, 'g') {
			return zc1771Hit(cmd, first, "global")
		}
		if strings.ContainsRune(first, 's') {
			return zc1771Hit(cmd, first, "suffix")
		}
	}
	return nil
}

func zc1771Hit(cmd *ast.SimpleCommand, flag, kind string) []Violation {
	return []Violation{{
		KataID: "ZC1771",
		Message: "`alias " + flag + "` defines a " + kind + " alias that expands outside " +
			"command position — a surprise for anyone reading the script later. Prefer a " +
			"function, or keep " + kind + " aliases in `~/.zshrc` where they are discoverable.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
