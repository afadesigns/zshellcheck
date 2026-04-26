// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1518",
		Title:    "Warn on `bash -p` — privileged mode (skips env sanitisation on setuid)",
		Severity: SeverityWarning,
		Description: "`bash -p` (and `-o privileged`) tells bash not to drop its effective UID/GID " +
			"and not to sanitize the environment when started on a setuid wrapper. It is " +
			"explicitly the flag you use to keep `BASH_ENV`, `SHELLOPTS`, and similar " +
			"attacker-controlled variables active while running as a more privileged user. " +
			"Almost no legitimate script needs `-p`; audit and remove.",
		Check: checkZC1518,
	})
}

func checkZC1518(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "bash" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-p" {
			return []Violation{{
				KataID: "ZC1518",
				Message: "`bash -p` keeps the privileged environment on a setuid wrapper — " +
					"almost never needed, audit and remove.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
