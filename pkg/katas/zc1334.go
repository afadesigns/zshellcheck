// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1334",
		Title:    "Avoid `type -p` — use `whence -p` in Zsh",
		Severity: SeverityWarning,
		Description: "`type -p` is a Bash flag that prints the path of a command. " +
			"Zsh `type` does not support `-p`. Use `whence -p` to get " +
			"the path of an external command in Zsh.",
		Check: checkZC1334,
		Fix:   fixZC1334,
	})
}

// fixZC1334 rewrites `type -p X` / `type -P X` to `whence -p X`. The
// span covers both the `type` command name and the `-p`/`-P` flag in a
// single edit — emitting the wider rewrite ensures it wins over the
// narrower `type` -> `command -v` swap from ZC1064 when both katas fire
// on the same input. Trailing argument(s) stay in place.
func fixZC1334(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "type" {
		return nil
	}
	var flag ast.Expression
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-p" || val == "-P" {
			flag = arg
			break
		}
	}
	if flag == nil {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || nameOff+len("type") > len(source) {
		return nil
	}
	if string(source[nameOff:nameOff+len("type")]) != "type" {
		return nil
	}
	flagTok := flag.TokenLiteralNode()
	flagOff := LineColToByteOffset(source, flagTok.Line, flagTok.Column)
	if flagOff < 0 || flagOff+2 > len(source) {
		return nil
	}
	if string(source[flagOff:flagOff+2]) != "-p" && string(source[flagOff:flagOff+2]) != "-P" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  flagOff + 2 - nameOff,
		Replace: "whence -p",
	}}
}

func checkZC1334(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "type" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-p" || val == "-P" {
			return []Violation{{
				KataID:  "ZC1334",
				Message: "Avoid `type -p` in Zsh — use `whence -p` to get the command path instead.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityWarning,
			}}
		}
	}

	return nil
}
