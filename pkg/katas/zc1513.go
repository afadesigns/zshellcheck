// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1513",
		Title:    "Style: `make install` without `DESTDIR=` — unmanaged system-wide install",
		Severity: SeverityStyle,
		Description: "`make install` drops files directly into `$(prefix)` with no package " +
			"manager tracking. Upgrades can leave stale files behind, uninstalls rely on " +
			"`make uninstall` being accurate, and the operation typically needs `sudo`. For " +
			"local builds, set `DESTDIR=/tmp/pkgroot` + wrap in `checkinstall` / `fpm` / " +
			"distro packaging, or use `stow` / `xstow` to manage symlinks under `/usr/local`.",
		Check: checkZC1513,
	})
}

func checkZC1513(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "make" && ident.Value != "gmake" {
		return nil
	}

	hasInstall := false
	hasDestdir := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "install" {
			hasInstall = true
		}
		if len(v) >= 8 && v[:8] == "DESTDIR=" {
			hasDestdir = true
		}
	}
	if !hasInstall || hasDestdir {
		return nil
	}
	return []Violation{{
		KataID: "ZC1513",
		Message: "`make install` without `DESTDIR=` leaves no package-manager record. Set " +
			"`DESTDIR=/tmp/pkgroot` and wrap in checkinstall / fpm, or use stow.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
