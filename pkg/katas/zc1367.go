// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1367",
		Title:    "Use Zsh `strftime` instead of Bash `printf '%(fmt)T'`",
		Severity: SeverityStyle,
		Description: "Bash 4.2+ supports `printf '%(fmt)T\\n' seconds` to format a timestamp. Zsh's " +
			"`zsh/datetime` module provides `strftime` which is more readable and works " +
			"consistently across versions: `strftime '%Y-%m-%d' $EPOCHSECONDS`.",
		Check: checkZC1367,
	})
}

func checkZC1367(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "printf" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		// Look for %(...)T format specifier
		if strings.Contains(val, ")T") && strings.Contains(val, "%(") {
			return []Violation{{
				KataID: "ZC1367",
				Message: "Use Zsh `strftime fmt seconds` (from `zsh/datetime`) instead of Bash " +
					"`printf '%(fmt)T' seconds`. Same formatting, more readable, no Bash-version gating.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
