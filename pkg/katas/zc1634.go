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
	if !ok || CommandIdentifier(cmd) != "umask" || len(cmd.Arguments) != 1 {
		return nil
	}
	v := cmd.Arguments[0].String()
	if !zc1634UmaskMissesWorldWrite(v) {
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

var zc1634AllZero = map[string]struct{}{"0": {}, "00": {}, "000": {}, "0000": {}}

func zc1634UmaskMissesWorldWrite(v string) bool {
	if _, hit := zc1634AllZero[v]; hit {
		return false
	}
	if len(v) < 3 || len(v) > 4 {
		return false
	}
	for _, c := range v {
		if c < '0' || c > '7' {
			return false
		}
	}
	switch v[len(v)-1] {
	case '0', '1', '4', '5':
		return true
	}
	return false
}
