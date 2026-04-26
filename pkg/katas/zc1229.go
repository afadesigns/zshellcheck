// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1229",
		Title:    "Prefer `rsync` over `scp` for file transfers",
		Severity: SeverityStyle,
		Description: "`scp` uses a deprecated protocol and lacks delta transfer, resume, " +
			"and progress features. `rsync` is more efficient and reliable for scripts.",
		Check: checkZC1229,
	})
}

func checkZC1229(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "scp" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1229",
		Message: "Prefer `rsync -az` over `scp` for file transfers. " +
			"`rsync` supports delta transfers, resume, and is more efficient.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
