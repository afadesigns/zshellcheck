// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1584",
		Title:    "Warn on `sudo -E` / `--preserve-env` — carries caller env into root shell",
		Severity: SeverityWarning,
		Description: "`sudo -E` preserves the invoking user's environment — `PATH`, " +
			"`LD_PRELOAD`, `PYTHONPATH`, etc. On a workstation where the user has a personal " +
			"`~/bin` early in `$PATH`, any wrapper named like a system binary gets executed " +
			"by the privileged process. That is exactly the sudoers `secure_path` mechanic " +
			"fails to protect against. Whitelist specific variables with `env_keep` in " +
			"sudoers, or call `sudo env VAR=value cmd` with the minimum.",
		Check: checkZC1584,
	})
}

func checkZC1584(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "sudo" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-E" || v == "--preserve-env" {
			return []Violation{{
				KataID: "ZC1584",
				Message: "`sudo " + v + "` carries the caller's PATH / LD_PRELOAD / … into " +
					"the privileged process. Use `env_keep` in sudoers or explicit `sudo " +
					"env VAR=… cmd`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
