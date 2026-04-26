// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1166",
		Title:    "Avoid `grep -i` for case-insensitive match — use `(#i)` glob flag",
		Severity: SeverityStyle,
		Description: "Zsh provides the `(#i)` glob flag for case-insensitive matching. " +
			"For variable matching, use `[[ $var == (#i)pattern ]]` instead of piping through grep -i.",
		Check: checkZC1166,
	})
}

func checkZC1166(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "grep" {
		return nil
	}

	hasCaseInsensitive := false
	hasFile := false
	patternSeen := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			if val == "-i" {
				hasCaseInsensitive = true
			}
		} else {
			if patternSeen {
				hasFile = true
				break
			}
			patternSeen = true
		}
	}

	if !hasCaseInsensitive || hasFile {
		return nil
	}

	return []Violation{{
		KataID: "ZC1166",
		Message: "Use Zsh `(#i)` glob flag for case-insensitive matching instead of piping through `grep -i`. " +
			"Example: `[[ $var == (#i)pattern ]]`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
