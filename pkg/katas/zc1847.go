// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1847",
		Title:    "Warn on `setopt CHASE_LINKS` — every `cd` silently swaps symlink paths for the real inode",
		Severity: SeverityWarning,
		Description: "`CHASE_LINKS` off is the Zsh default: `cd releases/current` leaves `$PWD` " +
			"as the logical path the user typed, and `cd ..` steps back up through the " +
			"symlink to where they came from. Turning the option on globally makes every " +
			"`cd` resolve the target to its physical inode — so `cd releases/current` lands " +
			"in `/srv/app/releases/20260415-deadbeef`, and the next `cd ../config` looks " +
			"for `/srv/app/releases/config` instead of the `/srv/app/config` that the user " +
			"expected. Scripts that rely on blue/green-style `current` symlinks break " +
			"silently. Keep the option off at the script level and request one-shot " +
			"physical resolution with `cd -P target` when a specific call needs it.",
		Check: checkZC1847,
	})
}

func checkZC1847(node ast.Node) []Violation {
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
			if zc1847IsChaseLinks(v) {
				return zc1847Hit(cmd, "setopt "+v)
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOCHASELINKS" {
				return zc1847Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1847IsChaseLinks(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "CHASELINKS"
}

func zc1847Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1847",
		Message: "`" + where + "` makes every `cd` resolve symlinks to the physical " +
			"inode — `cd releases/current` lands in the release dir, breaking `..` " +
			"navigation. Keep it off; use `cd -P target` one-shot when needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
