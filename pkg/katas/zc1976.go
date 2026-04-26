// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1976",
		Title:    "Error on `exportfs -au` / `exportfs -u` — unexports live NFS shares, clients get `ESTALE`",
		Severity: SeverityError,
		Description: "`exportfs -au` unexports every NFS share on the server; `exportfs -u " +
			"HOST:/PATH` removes a single share. Any client that currently has the " +
			"export mounted is not notified — the next read/write returns `ESTALE`, " +
			"the mount looks live but every open fd fails, and the only recovery is a " +
			"client-side `umount -l` + remount. `exportfs -f` (flush) is almost always " +
			"what you actually want after an `/etc/exports` edit; keep `-u`/`-au` for " +
			"planned shutdowns with a coordinated client `umount` first.",
		Check: checkZC1976,
	})
}

func checkZC1976(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "exportfs" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-au" || v == "-ua" || v == "-u" {
			return zc1976Hit(cmd, "exportfs "+v)
		}
	}
	return nil
}

func zc1976Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1976",
		Message: "`" + form + "` unexports live NFS shares — mounted clients see " +
			"`ESTALE` on every open fd. Use `exportfs -f` after editing " +
			"`/etc/exports`; reserve `-u`/`-au` for coordinated shutdowns.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
