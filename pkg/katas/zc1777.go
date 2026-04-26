// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

const zc1777PreloadPath = "/etc/ld.so.preload"

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1777",
		Title:    "Error on `tee/cp/mv/install/dd` writing `/etc/ld.so.preload` — classic rootkit persistence",
		Severity: SeverityError,
		Description: "`/etc/ld.so.preload` lists shared libraries that the dynamic linker " +
			"forcibly loads into every dynamically-linked binary, root processes included. " +
			"The file is almost never needed on a modern distribution — package managers " +
			"do not touch it, and `LD_PRELOAD` handles the per-invocation case without " +
			"persisting the change. A script that pipes content into `/etc/ld.so.preload` " +
			"with `tee` / `cp` / `mv` / `install` / `dd` is a textbook rootkit persistence " +
			"primitive (`libprocesshider`, `Azazel`, `Jynx`). Remove the line, audit " +
			"`/etc/ld.so.preload` for unexpected entries (`sha256sum`, `diff` against a " +
			"known-good backup), and if preloading is legitimately required, use a scoped " +
			"`LD_PRELOAD=` on the specific invocation.",
		Check: checkZC1777,
	})
}

func checkZC1777(node ast.Node) []Violation {
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
			if arg.String() == zc1777PreloadPath {
				return zc1777Hit(cmd, ident.Value+" "+zc1777PreloadPath)
			}
		}
	}

	prevRedir := ""
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevRedir != "" {
			if v == zc1777PreloadPath {
				return zc1777Hit(cmd, prevRedir+" "+zc1777PreloadPath)
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

func zc1777Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1777",
		Message: "`" + what + "` writes `/etc/ld.so.preload` — linker force-loads " +
			"each listed library into every process. Audit for unexpected entries; " +
			"for a scoped preload use `LD_PRELOAD=` on a single invocation.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
