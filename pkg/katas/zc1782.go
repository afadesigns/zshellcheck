// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1782",
		Title:    "Error on `flatpak remote-add --no-gpg-verify` — trust chain disabled for the repo",
		Severity: SeverityError,
		Description: "A Flatpak remote without GPG verification accepts any OSTree update that " +
			"the server (or anyone on the path) cares to send. Signatures are what connect " +
			"`flatpak install FOO` to the operator that actually built `FOO` — strip them and " +
			"the install reduces to a plain HTTPS download with no identity attached. If you " +
			"genuinely need a local / air-gapped repo, sign it yourself with `ostree gpg-sign` " +
			"and add the key via `--gpg-import=KEYFILE`. Never leave `--no-gpg-verify` in " +
			"provisioning scripts for production systems.",
		Check: checkZC1782,
	})
}

func checkZC1782(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "flatpak" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "remote-add" && cmd.Arguments[0].String() != "remote-modify" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--no-gpg-verify" ||
			v == "--gpg-verify=false" ||
			v == "--no-gpg-verify=true" {
			return []Violation{{
				KataID: "ZC1782",
				Message: "`flatpak " + cmd.Arguments[0].String() + " " + v + "` disables " +
					"signature verification — updates from this remote are accepted with " +
					"only HTTPS as identity. Sign the repo (`ostree gpg-sign`) and import " +
					"the key with `--gpg-import=KEYFILE`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
