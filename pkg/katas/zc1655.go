// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strconv"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1655",
		Title:    "Warn on `read -n N` — Bash reads N chars; Zsh's `-n` means \"drop newline\"",
		Severity: SeverityWarning,
		Description: "In Bash, `read -n N var` reads exactly N characters (handy for single-" +
			"keypress prompts). In Zsh, `-n` is the \"don't append newline to the reply " +
			"string\" flag and doesn't take a count — `read -n 1 var` sets `var` to the " +
			"whole line, not a single character. Use `read -k N var` in Zsh for N-character " +
			"reads.",
		Check: checkZC1655,
	})
}

func checkZC1655(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "read" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-n" && i+1 < len(cmd.Arguments) {
			if _, err := strconv.Atoi(cmd.Arguments[i+1].String()); err == nil {
				return zc1655Hit(cmd)
			}
		}
		if strings.HasPrefix(v, "-n") && len(v) > 2 {
			if _, err := strconv.Atoi(v[2:]); err == nil {
				return zc1655Hit(cmd)
			}
		}
	}
	return nil
}

func zc1655Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1655",
		Message: "`read -n N` is Bash syntax for \"read N characters\". Zsh's `-n` means " +
			"\"drop trailing newline\" with no count. Use `read -k N var` in Zsh.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
