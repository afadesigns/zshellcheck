// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1285",
		Title:    "Use Zsh `${(o)array}` for sorting instead of piping to `sort`",
		Severity: SeverityStyle,
		Description: "Zsh provides the `(o)` parameter expansion flag to sort array elements " +
			"in ascending order and `(O)` for descending order. This avoids spawning " +
			"an external `sort` process for simple array sorting.",
		Check: checkZC1285,
	})
}

func checkZC1285(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sort" {
		return nil
	}

	// sort with complex flags has legitimate uses beyond simple array sorting
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-t" || val == "-k" || val == "-n" || val == "-r" ||
			val == "-u" || val == "-h" || val == "-V" || val == "-g" ||
			val == "-c" || val == "-m" || val == "-s" {
			return nil
		}
	}

	// sort with only a filename — suggest Zsh native sorting
	if len(cmd.Arguments) == 1 {
		return []Violation{{
			KataID:  "ZC1285",
			Message: "Use Zsh `${(o)array}` for sorting instead of piping to `sort`. The `(o)` flag sorts in-shell.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}
