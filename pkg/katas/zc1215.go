// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1215",
		Title:    "Source `/etc/os-release` instead of parsing with `cat`/`grep`",
		Severity: SeverityStyle,
		Description: "`/etc/os-release` is designed to be sourced directly. " +
			"Use `. /etc/os-release` to get variables like `$ID`, `$VERSION_ID` without parsing.",
		Check: checkZC1215,
		Fix:   fixZC1215,
	})
}

// fixZC1215 rewrites `cat /etc/os-release` (or `/etc/lsb-release`) to
// `. /etc/os-release`. Single-edit replacement of the `cat` command
// name with the source builtin `.`. Only fires when cat has exactly
// one argument; piped or multi-file shapes are left alone. Idempotent
// — a re-run sees `.`, not `cat`. Defensive byte-match guard.
func fixZC1215(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cat" {
		return nil
	}
	if len(cmd.Arguments) != 1 {
		return nil
	}
	val := cmd.Arguments[0].String()
	if val != "/etc/os-release" && val != "/etc/lsb-release" {
		return nil
	}
	off := LineColToByteOffset(source, v.Line, v.Column)
	if off < 0 || off+len("cat") > len(source) {
		return nil
	}
	if string(source[off:off+len("cat")]) != "cat" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("cat"),
		Replace: ".",
	}}
}

func checkZC1215(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cat" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "/etc/os-release" || val == "/etc/lsb-release" {
			return []Violation{{
				KataID: "ZC1215",
				Message: "Source `/etc/os-release` directly with `. /etc/os-release` instead of " +
					"parsing with `cat`. It exports variables like `$ID` and `$VERSION_ID`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
