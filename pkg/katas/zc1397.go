// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1397",
		Title:    "Avoid `$COMP_TYPE`/`$COMP_KEY` — Bash completion globals, not in Zsh",
		Severity: SeverityError,
		Description: "Bash programmable completion exposes `$COMP_TYPE` (completion type) and " +
			"`$COMP_KEY` (completion key pressed). Zsh's compsys does not use these variables; " +
			"query completion context via `$compstate` assoc array or context keys from " +
			"`_arguments`/`_values` instead.",
		Check: checkZC1397,
	})
}

func checkZC1397(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "echo" && ident.Value != "print" && ident.Value != "printf" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "COMP_TYPE") || strings.Contains(v, "COMP_KEY") ||
			strings.Contains(v, "COMP_WORDBREAKS") {
			return []Violation{{
				KataID: "ZC1397",
				Message: "Bash `$COMP_TYPE`/`$COMP_KEY`/`$COMP_WORDBREAKS` are not Zsh-native. " +
					"Use `$compstate` associative array for completion context in Zsh.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}
