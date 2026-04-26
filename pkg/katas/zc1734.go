// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1734IdentityFiles = map[string]bool{
	"/etc/passwd":  true,
	"/etc/shadow":  true,
	"/etc/group":   true,
	"/etc/gshadow": true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1734",
		Title:    "Error on `cp/mv/tee` overwriting `/etc/passwd|shadow|group|gshadow`",
		Severity: SeverityError,
		Description: "The user-identity files are managed by `useradd` / `usermod` / `vipw` / " +
			"`vigr`, which take a file lock and keep `passwd` / `shadow` (and `group` / " +
			"`gshadow`) in sync. Replacing them with `cp`, `mv`, `tee`, or a redirect " +
			"(`echo … > /etc/passwd`) bypasses the lock: concurrent edits race, malformed " +
			"entries lock the whole system out, and the shadow file ends up pointing at " +
			"users that no longer exist. Use `vipw -e` / `vigr -e` to edit, or `useradd` " +
			"/ `usermod` / `passwd` to mutate one entry at a time.",
		Check: checkZC1734,
	})
}

func checkZC1734(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "cp", "mv", "tee", "install", "dd":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if zc1734IdentityFiles[v] {
				return zc1734Hit(cmd, ident.Value+" "+v)
			}
		}
	}

	// Redirect form: any command whose args contain `>` or `>>` followed by an
	// identity file path.
	prevRedir := ""
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevRedir != "" {
			if zc1734IdentityFiles[v] {
				return zc1734Hit(cmd, prevRedir+" "+v)
			}
			prevRedir = ""
			continue
		}
		if v == ">" || v == ">>" {
			prevRedir = v
		}
	}
	return nil
}

func zc1734Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1734",
		Message: "`" + what + "` bypasses the lock that `vipw` / `vigr` / `useradd` use " +
			"on the user-identity files. Edit through `vipw -e` / `vigr -e`, or mutate " +
			"a single entry with `useradd` / `usermod` / `passwd`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
