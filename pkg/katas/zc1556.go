// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1556",
		Title:    "Error on `openssl enc -des` / `-rc4` / `-3des` — broken symmetric cipher",
		Severity: SeverityError,
		Description: "DES, RC4, and 3DES are all broken or on-deprecation-path: DES's 56-bit key " +
			"fell to commodity brute-force decades ago, RC4 has practical biased-output attacks, " +
			"and 3DES suffers the Sweet32 birthday collision when reused for more than ~32GB. " +
			"None of them provide authenticity either. Use `-aes-256-gcm` or `-chacha20-poly1305`, " +
			"or move up to a dedicated tool (`age`, `gpg`, `libsodium`).",
		Check: checkZC1556,
	})
}

func checkZC1556(node ast.Node) []Violation {
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

	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "enc" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := strings.ToLower(arg.String())
		switch v {
		case "-des", "-des-cbc", "-des-ecb",
			"-des3", "-des-ede", "-des-ede-cbc", "-des-ede3", "-des-ede3-cbc",
			"-rc4", "-rc4-40",
			"-bf", "-bf-cbc",
			"-rc2", "-rc2-cbc",
			"-cast", "-cast5-cbc":
			return []Violation{{
				KataID: "ZC1556",
				Message: "`openssl enc " + v + "` is a broken or deprecated cipher. Use " +
					"`-aes-256-gcm` / `-chacha20-poly1305`, or `age` / `gpg` for files.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
