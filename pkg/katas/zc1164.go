// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1164",
		Title:    "Avoid `sed -n 'Np'` — use Zsh array subscript",
		Severity: SeverityStyle,
		Description: "Extracting a specific line with `sed -n 'Np'` spawns a process. " +
			"Use Zsh array subscript `${lines[N]}` after splitting with `${(f)...}`.",
		Check: checkZC1164,
	})
}

func checkZC1164(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sed" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}

	first := cmd.Arguments[0].String()
	if first != "-n" {
		return nil
	}

	// Check if second arg matches pattern like '3p', '10p', etc.
	second := strings.Trim(cmd.Arguments[1].String(), "'\"")
	if len(second) >= 2 && second[len(second)-1] == 'p' {
		allDigits := true
		for _, ch := range second[:len(second)-1] {
			if ch < '0' || ch > '9' {
				allDigits = false
				break
			}
		}
		if allDigits && len(cmd.Arguments) == 2 {
			return []Violation{{
				KataID: "ZC1164",
				Message: "Use Zsh array subscript `${lines[N]}` instead of `sed -n 'Np'`. " +
					"Split input with `${(f)...}` then index directly.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
