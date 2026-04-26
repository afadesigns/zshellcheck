// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1519",
		Title:    "Warn on `ulimit -u unlimited` — removes user process cap, enables fork bombs",
		Severity: SeverityWarning,
		Description: "`ulimit -u` caps the number of processes a UID can run; `unlimited` removes " +
			"that cap. Combined with a bug in a background loop (or a literal fork bomb via " +
			"`:(){ :|:& };:`) it pegs the scheduler until the machine has to be cold-booted. " +
			"Pick a realistic number (distro defaults around 4096 for interactive sessions) or " +
			"set it in `/etc/security/limits.d/` so it is persistent and visible.",
		Check: checkZC1519,
	})
}

func checkZC1519(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ulimit" {
		return nil
	}

	var prevU bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevU {
			prevU = false
			if v == "unlimited" {
				return []Violation{{
					KataID: "ZC1519",
					Message: "`ulimit -u unlimited` removes the user process cap — fork bomb " +
						"surface. Pick a realistic number or set it via /etc/security/limits.d/.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
		if v == "-u" {
			prevU = true
		}
	}
	return nil
}
