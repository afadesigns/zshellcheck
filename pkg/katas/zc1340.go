// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1340",
		Title:    "Avoid `shuf` for random array element — use Zsh `$RANDOM`",
		Severity: SeverityStyle,
		Description: "Zsh provides `$RANDOM` and array subscripts to pick random elements " +
			"without spawning `shuf`. For a single random array element, use " +
			"`${array[RANDOM%$#array+1]}`.",
		Check: checkZC1340,
	})
}

func checkZC1340(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "shuf" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1340",
		Message: "Avoid `shuf` for random selection — use Zsh `${array[RANDOM%$#array+1]}` " +
			"with `$RANDOM` for in-shell randomness without spawning an external.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
