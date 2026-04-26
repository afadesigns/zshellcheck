// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1570",
		Title:    "Warn on `smbclient -N` / `mount.cifs guest` — anonymous SMB share access",
		Severity: SeverityWarning,
		Description: "`smbclient -N` skips authentication entirely (anonymous / null session); " +
			"`mount.cifs` with `guest,username=` or `-o guest` does the same at the mount " +
			"layer. Any host on the network segment can then read the share. If the share is " +
			"truly public (software mirror, build cache) wrap in a read-only filesystem and " +
			"document it; otherwise require Kerberos (`-k`) or pass credentials via " +
			"`credentials=<file>` with 0600 perms.",
		Check: checkZC1570,
	})
}

func checkZC1570(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "smbclient" && ident.Value != "mount.cifs" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-N" || v == "--no-pass" {
			return []Violation{{
				KataID: "ZC1570",
				Message: "`" + ident.Value + " " + v + "` is anonymous SMB access — any " +
					"host on-net can read the share. Use credentials=<file> 0600 or -k.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
