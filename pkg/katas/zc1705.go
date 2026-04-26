// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1705",
		Title:    "Info: `awk -i inplace` is gawk-only — script breaks on mawk / BSD awk",
		Severity: SeverityInfo,
		Description: "The `inplace` extension that powers `awk -i inplace` ships only with gawk. " +
			"On Alpine (default `mawk`), Debian-busybox, macOS, FreeBSD, NetBSD, OpenBSD, " +
			"or any container image without `gawk` installed the script aborts with " +
			"`fatal: can't open extension 'inplace'`. If portability matters, write through " +
			"a temporary file (`awk … input > tmp && mv tmp input`); if you really do need " +
			"in-place edits in scripts that target gawk only, document the requirement and " +
			"add `command -v gawk >/dev/null` at the top.",
		Check: checkZC1705,
	})
}

func checkZC1705(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "awk" && ident.Value != "gawk" && ident.Value != "mawk" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		if arg.String() != "-i" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		if cmd.Arguments[i+1].String() == "inplace" {
			return []Violation{{
				KataID: "ZC1705",
				Message: "`awk -i inplace` is gawk-only — fails on mawk / BSD awk / busybox " +
					"awk. For portability rewrite as `awk … input > tmp && mv tmp input`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}
	return nil
}
