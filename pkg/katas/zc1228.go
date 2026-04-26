// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1228",
		Title:    "Avoid `ssh` without host key policy in scripts",
		Severity: SeverityWarning,
		Description: "`ssh` without `-o BatchMode=yes` or `-o StrictHostKeyChecking` prompts " +
			"interactively for host key verification, hanging non-interactive scripts.",
		Check: checkZC1228,
	})
}

func checkZC1228(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ssh" {
		return nil
	}

	hasBatchOrStrict := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "BatchMode=yes" || val == "StrictHostKeyChecking=accept-new" ||
			val == "StrictHostKeyChecking=no" || val == "StrictHostKeyChecking=yes" {
			hasBatchOrStrict = true
		}
	}

	if !hasBatchOrStrict {
		return []Violation{{
			KataID: "ZC1228",
			Message: "Use `ssh -o BatchMode=yes` or `-o StrictHostKeyChecking=accept-new` in scripts. " +
				"Without these, ssh may prompt interactively and hang.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
