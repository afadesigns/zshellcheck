// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1474",
		Title:    "Warn on `ssh-keygen -N \"\"` — generates passwordless SSH key",
		Severity: SeverityWarning,
		Description: "Generating an SSH key with an empty passphrase (`-N \"\"`) leaves the key " +
			"usable by anything that can read the file. Combined with a weak umask or a backup " +
			"that follows the file, this is a common lateral-movement vector. Use a real " +
			"passphrase, or delegate key storage to `ssh-agent` / a hardware token.",
		Check: checkZC1474,
	})
}

func checkZC1474(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ssh-keygen" {
		return nil
	}

	var prevN bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevN {
			prevN = false
			if v == `""` || v == `''` || v == "" {
				return []Violation{{
					KataID:  "ZC1474",
					Message: "`ssh-keygen -N \"\"` generates a passwordless key — anything that reads the file can use it. Use a passphrase or ssh-agent/HSM.",
					Line:    cmd.Token.Line,
					Column:  cmd.Token.Column,
					Level:   SeverityWarning,
				}}
			}
		}
		if v == "-N" {
			prevN = true
		}
	}
	return nil
}
