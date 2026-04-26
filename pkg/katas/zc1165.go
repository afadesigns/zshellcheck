// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1165",
		Title:    "Use Zsh parameter expansion for simple `awk` field extraction",
		Severity: SeverityStyle,
		Description: "Simple `awk '{print $1}'` or `awk '{print $NF}'` can often be replaced with " +
			"Zsh parameter expansion `${var%% *}` (first field) or `${var##* }` (last field).",
		Check: checkZC1165,
	})
}

func checkZC1165(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "awk" {
		return nil
	}

	// Only flag awk with a single print statement and no file argument
	if len(cmd.Arguments) != 1 {
		return nil
	}

	arg := strings.Trim(cmd.Arguments[0].String(), "'\"")
	if arg == "{print $1}" || arg == "{print $NF}" {
		return []Violation{{
			KataID: "ZC1165",
			Message: "Use Zsh parameter expansion (`${var%% *}` or `${var##* }`) instead of " +
				"`awk '{print $1}'` for simple field extraction without spawning awk.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
