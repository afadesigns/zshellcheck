// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1569",
		Title:    "Error on `nvme format -s1` / `-s2` — cryptographic or full-block SSD erase",
		Severity: SeverityError,
		Description: "`nvme format -s1` does a cryptographic erase of the target namespace; " +
			"`-s2` (or the full-NVMe sanitize) rewrites every block. Both are unrecoverable " +
			"in seconds. On a typo in the device variable — or a script that iterates over " +
			"`/dev/nvme*n*` and catches the wrong namespace — the wrong disk is gone by the " +
			"time the operator notices. Run interactively on verified targets, or not at all " +
			"from automation.",
		Check: checkZC1569,
	})
}

func checkZC1569(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "nvme" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "format" && sub != "sanitize" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "-s1" || v == "-s2" || v == "--ses=1" || v == "--ses=2" ||
			v == "-a" || v == "--sanact" {
			return []Violation{{
				KataID: "ZC1569",
				Message: "`nvme " + sub + " " + v + "` unrecoverably erases the namespace in " +
					"seconds. Do not run from automation.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
