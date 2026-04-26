// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1686",
		Title:    "Warn on `compinit -C` / `compinit -u` — skips / ignores `$fpath` integrity checks",
		Severity: SeverityWarning,
		Description: "Zsh's completion system loads every file from `$fpath` as shell code. " +
			"`compinit` normally warns when an `$fpath` directory (or a file in one) is " +
			"writable by someone other than the current user or root, and skips loading. " +
			"`compinit -C` skips the security check entirely for speed; `compinit -u` " +
			"acknowledges the warning and loads the insecure files anyway. Either way, a " +
			"world-writable entry in `$fpath` becomes an execution primitive for any user " +
			"on the host. Audit `$fpath` with `compaudit`, fix ownership / permissions, " +
			"then run plain `compinit`.",
		Check: checkZC1686,
	})
}

func checkZC1686(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "compinit" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-C" {
			return zc1686Hit(cmd, "-C", "skip-security-check")
		}
		if v == "-u" {
			return zc1686Hit(cmd, "-u", "load-insecure-files")
		}
	}
	return nil
}

func zc1686Hit(cmd *ast.SimpleCommand, flag, what string) []Violation {
	return []Violation{{
		KataID: "ZC1686",
		Message: "`compinit " + flag + "` (" + what + ") loads `$fpath` files that are " +
			"writable by others — any user on the host can inject shell code. Run " +
			"`compaudit`, fix permissions, then `compinit` without the flag.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
