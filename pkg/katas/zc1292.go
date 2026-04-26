// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1292",
		Title:    "Use Zsh `${var//old/new}` instead of `tr` for character translation",
		Severity: SeverityStyle,
		Description: "Zsh provides `${var//old/new}` for global substitution within a variable. " +
			"For simple single-character translation, this avoids spawning `tr` as an " +
			"external process.",
		Check: checkZC1292,
	})
}

func checkZC1292(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "tr" {
		return nil
	}

	if len(cmd.Arguments) != 2 {
		return nil
	}

	first := strings.Trim(cmd.Arguments[0].String(), "'\"")
	second := strings.Trim(cmd.Arguments[1].String(), "'\"")

	// Only flag simple single-character translations, not ranges or classes
	if len(first) == 1 && len(second) == 1 {
		return []Violation{{
			KataID:  "ZC1292",
			Message: "Use Zsh `${var//" + first + "/" + second + "}` for character substitution instead of `tr`.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}
