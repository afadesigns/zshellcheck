// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1707",
		Title:    "Warn on `gpg --keyserver hkp://…` — plaintext keyserver fetch",
		Severity: SeverityWarning,
		Description: "`hkp://` is the unencrypted HKP keyserver protocol. A MITM on the path " +
			"(corporate proxy, hotel Wi-Fi, hostile router) can swap key bytes during the " +
			"fetch and `gpg --recv-keys` happily imports the substitute. Use `hkps://" +
			"keys.openpgp.org` (TLS) or fetch the armored key over HTTPS and verify the " +
			"fingerprint out-of-band before `gpg --import`.",
		Check: checkZC1707,
	})
}

func checkZC1707(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value == "keyserver" {
		// Parser-mangled form: `gpg --keyserver hkp://…` lost `gpg`.
		if len(cmd.Arguments) > 0 && strings.HasPrefix(cmd.Arguments[0].String(), "hkp://") {
			return zc1707Hit(cmd)
		}
		return nil
	}
	if ident.Value != "gpg" && ident.Value != "gpg2" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "--keyserver=hkp://") {
			return zc1707Hit(cmd)
		}
		if v == "--keyserver" && i+1 < len(cmd.Arguments) &&
			strings.HasPrefix(cmd.Arguments[i+1].String(), "hkp://") {
			return zc1707Hit(cmd)
		}
	}
	return nil
}

func zc1707Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1707",
		Message: "`gpg --keyserver hkp://…` is plaintext — a MITM swaps the key bytes. Use " +
			"`hkps://keys.openpgp.org` or fetch over HTTPS and verify the fingerprint.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
