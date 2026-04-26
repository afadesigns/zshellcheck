// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1990",
		Title:    "Warn on `openssl passwd -crypt` / `-1` / `-apr1` — obsolete password hash formats",
		Severity: SeverityWarning,
		Description: "`openssl passwd -crypt` emits DES-crypt, 8-char truncated and crackable in " +
			"seconds on modern hardware. `-1` is FreeBSD-style MD5, unsuitable for " +
			"storage, long broken. `-apr1` is Apache's MD5-based variant with the same " +
			"weakness. Any hash produced by these flags lands in `/etc/shadow`, an " +
			"htpasswd file, or a database row where an attacker can offline-crack the " +
			"whole batch with a single GPU. Use `-5` (SHA-256-crypt), `-6` (SHA-512-" +
			"crypt), or prefer a dedicated KDF-based hasher — `mkpasswd -m yescrypt`, " +
			"`htpasswd -B` (bcrypt), or `argon2` — so brute-force cost scales with " +
			"hardware.",
		Check: checkZC1990,
	})
}

func checkZC1990(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "openssl" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "passwd" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		switch v {
		case "-crypt", "-1", "-apr1":
			return []Violation{{
				KataID: "ZC1990",
				Message: "`openssl passwd " + v + "` emits a broken hash format — " +
					"DES/MD5 variants crack on a laptop. Use `-5` / `-6` or a " +
					"KDF-based hasher (`mkpasswd -m yescrypt`, `htpasswd -B`, " +
					"`argon2`).",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
