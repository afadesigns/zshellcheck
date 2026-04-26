// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1477",
		Title:    "Warn on `printf \"$var\"` — variable in format-string position (printf-fmt attack)",
		Severity: SeverityWarning,
		Description: "The first argument to `printf` is a format string. Interpolating a shell " +
			"variable into it means any `%` sequence inside the variable is interpreted as a " +
			"format specifier — at best producing garbage, at worst crashing with " +
			"`%s`-out-of-bounds reads or writing attacker-controlled data with `%n`. Always " +
			"use a literal format string: `printf '%s\\n' \"$var\"`.",
		Check: checkZC1477,
	})
}

func checkZC1477(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "printf" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}
	first := cmd.Arguments[0].String()
	raw := stripOuterQuotes(first)

	// Single-quoted strings don't interpolate; treat them as safe even if `$` is present.
	if strings.HasPrefix(first, "'") && strings.HasSuffix(first, "'") {
		return nil
	}

	// Look for an unescaped `$` (variable, command substitution, or arithmetic).
	for i := 0; i < len(raw); i++ {
		if raw[i] == '\\' {
			i++
			continue
		}
		if raw[i] == '$' {
			return []Violation{{
				KataID: "ZC1477",
				Message: "`printf` format string contains a variable — `%` inside `$var` is " +
					"reparsed as a format specifier. Use `printf '%s' \"$var\"` instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
