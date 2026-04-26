// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1185",
		Title:    "Use Zsh `${#${(z)var}}` instead of `wc -w` for word count",
		Severity: SeverityStyle,
		Description: "Zsh `${(z)var}` splits a string into words and `${#...}` counts them. " +
			"Avoid piping through `wc -w` for simple word counting from variables.",
		Check: checkZC1185,
	})
}

func checkZC1185(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "wc" {
		return nil
	}

	hasWordFlag := false
	hasFile := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-w" {
			hasWordFlag = true
		} else if len(val) > 0 && val[0] != '-' {
			hasFile = true
		}
	}

	if hasWordFlag && !hasFile {
		return []Violation{{
			KataID: "ZC1185",
			Message: "Use Zsh `${#${(z)var}}` for word counting instead of piping through `wc -w`. " +
				"Parameter expansion avoids spawning an external process.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
