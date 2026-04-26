// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1549",
		Title:    "Error on `unzip -d /` / `unzip -o ... -d /` — extract archive into filesystem root",
		Severity: SeverityError,
		Description: "Unzipping directly into `/` (or `/root`, `/boot`) overwrites any system file " +
			"whose path matches an entry in the archive. A malicious zip that carries " +
			"`etc/passwd`, `usr/bin/ls`, or `root/.ssh/authorized_keys` turns a seemingly " +
			"harmless extract into full system compromise. Stage to a scratch directory, " +
			"inspect contents, then copy or install specific files.",
		Check: checkZC1549,
	})
}

func checkZC1549(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "unzip" && ident.Value != "busybox" {
		return nil
	}

	var prevD bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevD {
			prevD = false
			if v == "/" || v == "/root" || v == "/boot" {
				return []Violation{{
					KataID: "ZC1549",
					Message: "`unzip -d " + v + "` extracts into a system path — any archive " +
						"entry overwrites matching system file. Stage, inspect, copy.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityError,
				}}
			}
		}
		if v == "-d" {
			prevD = true
		}
	}
	return nil
}
