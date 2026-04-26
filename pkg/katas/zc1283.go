// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1283",
		Title:    "Use `setopt` instead of `set -o` for Zsh options",
		Severity: SeverityStyle,
		Description: "Zsh provides `setopt` and `unsetopt` as native builtins for managing shell " +
			"options. Using `set -o` / `set +o` is a POSIX compatibility form that is less " +
			"idiomatic in Zsh scripts.",
		Check: checkZC1283,
		Fix:   fixZC1283,
	})
}

// fixZC1283 rewrites `set -o OPTION` into `setopt OPTION`. The span
// covers the `set` command name and the `-o` flag in a single edit;
// trailing option arguments stay in place.
func fixZC1283(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "set" {
		return nil
	}
	var dashO ast.Expression
	for _, arg := range cmd.Arguments {
		if arg.String() == "-o" {
			dashO = arg
			break
		}
	}
	if dashO == nil {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || nameOff+len("set") > len(source) {
		return nil
	}
	if string(source[nameOff:nameOff+len("set")]) != "set" {
		return nil
	}
	dashTok := dashO.TokenLiteralNode()
	dashOff := LineColToByteOffset(source, dashTok.Line, dashTok.Column)
	if dashOff < 0 || dashOff+2 > len(source) {
		return nil
	}
	if string(source[dashOff:dashOff+2]) != "-o" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  dashOff + 2 - nameOff,
		Replace: "setopt",
	}}
}

func checkZC1283(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "set" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-o" {
			return []Violation{{
				KataID:  "ZC1283",
				Message: "Use `setopt` instead of `set -o` in Zsh scripts. `setopt` is the native Zsh idiom.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityStyle,
			}}
		}
	}

	return nil
}
