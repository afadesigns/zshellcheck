// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1769",
		Title:    "Warn on `vagrant destroy --force` — VM destroyed without confirmation",
		Severity: SeverityWarning,
		Description: "`vagrant destroy --force` (alias `-f`) tears every VM in the Vagrantfile " +
			"down — and their ephemeral filesystem state — without prompting. Any data " +
			"provisioned into the VM that was never exported back to the host (database " +
			"seeds, build caches, local-only test fixtures) goes with it. In unattended " +
			"scripts, drop the flag so the prompt still gates the destroy; for CI cycles, " +
			"`vagrant halt` + `vagrant up` reuses the same box without losing state.",
		Check: checkZC1769,
	})
}

func checkZC1769(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "vagrant" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "destroy" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--force" || v == "-f" {
			return []Violation{{
				KataID: "ZC1769",
				Message: "`vagrant destroy " + v + "` skips the prompt and drops the VM " +
					"(and any un-exported data inside). Drop the flag, or use `vagrant " +
					"halt` + `vagrant up` to cycle without destroy.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
