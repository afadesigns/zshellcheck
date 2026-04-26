// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1476",
		Title:    "Warn on `apt-key add` — deprecated, trusts every repo system-wide",
		Severity: SeverityWarning,
		Description: "`apt-key` was deprecated in APT 2.2 and removed from `apt` 2.5. Keys added " +
			"with `apt-key add` end up in a global keyring that signs every repo on the system, " +
			"so a typo-squatted third-party PPA can ship updates for `apt`, `libc6`, or " +
			"`openssh-server`. Store the key in `/etc/apt/keyrings/<vendor>.gpg` and scope it in " +
			"`signed-by=` on the specific `sources.list.d` entry.",
		Check: checkZC1476,
	})
}

func checkZC1476(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "apt-key" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub == "add" || sub == "adv" {
		return []Violation{{
			KataID: "ZC1476",
			Message: "`apt-key " + sub + "` adds to a global keyring that signs every repo. " +
				"Use `/etc/apt/keyrings/<vendor>.gpg` + `signed-by=` instead.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}
