// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1364",
		Title:    "Use Zsh `${var:pos:len}` instead of `cut -c` for character ranges",
		Severity: SeverityStyle,
		Description: "`cut -c N-M` extracts characters N through M from each line. Zsh's " +
			"`${var:pos:len}` (0-indexed position, length) does the same from a variable " +
			"without spawning `cut`.",
		Check: checkZC1364,
	})
}

func checkZC1364(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cut" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := strings.TrimFunc(arg.String(), func(r rune) bool { return r == '\'' || r == '"' })
		if val == "-c" || val == "--characters" ||
			(len(val) > 2 && val[:2] == "-c") ||
			strings.HasPrefix(val, "--characters=") {
			return []Violation{{
				KataID: "ZC1364",
				Message: "Use Zsh `${var:pos:len}` for character ranges instead of `cut -c`. " +
					"Parameter expansion is in-shell and zero-indexed.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
