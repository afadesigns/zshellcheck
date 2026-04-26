// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1634",
		Title:    "Warn on `umask NNN` that fails to mask world-write — mask-inversion footgun",
		Severity: SeverityWarning,
		Description: "`umask` is a mask: bits that are set are removed from the default " +
			"permission. The classic pitfall is reading it as \"permissions I want\" — " +
			"`umask 111` feels tight (\"no execute for anyone\") but it does not mask the write " +
			"bit, so every new file is `666` (rw-rw-rw-). The \"other\" digit must be one of " +
			"`2/3/6/7` to strip world-write. Use `022` for publicly readable files, `077` for " +
			"secrets-handling.",
		Check: checkZC1634,
	})
}

func checkZC1634(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "umask" {
		return nil
	}
	if len(cmd.Arguments) != 1 {
		return nil
	}

	v := cmd.Arguments[0].String()
	// Exclude forms already flagged by ZC1195 / ZC1516.
	if v == "0" || v == "00" || v == "000" || v == "0000" {
		return nil
	}
	if len(v) < 3 || len(v) > 4 {
		return nil
	}
	for _, c := range v {
		if c < '0' || c > '7' {
			return nil
		}
	}
	last := v[len(v)-1]
	if last != '0' && last != '1' && last != '4' && last != '5' {
		return nil
	}
	return []Violation{{
		KataID: "ZC1634",
		Message: "`umask " + v + "` leaves world-write on new files — the \"other\" digit " +
			"must be `2`/`3`/`6`/`7` to mask the write bit. Use `022` for public, `077` " +
			"for secrets.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
