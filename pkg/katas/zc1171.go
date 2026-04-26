// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1171",
		Title:    "Use `print` instead of `echo -e` for escape sequences",
		Severity: SeverityStyle,
		Description: "`echo -e` behavior varies across shells and platforms. " +
			"In Zsh, `print` natively interprets escape sequences and is more reliable.",
		Check: checkZC1171,
		Fix:   fixZC1171,
	})
}

// fixZC1171 collapses `echo -e` into `print`. Span covers the
// command name, intervening whitespace, and the `-e` flag; remaining
// arguments stay in place.
func fixZC1171(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "echo" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || nameOff+len("echo") > len(source) {
		return nil
	}
	if string(source[nameOff:nameOff+len("echo")]) != "echo" {
		return nil
	}
	i := nameOff + len("echo")
	for i < len(source) && (source[i] == ' ' || source[i] == '\t') {
		i++
	}
	if i+2 > len(source) || source[i] != '-' || source[i+1] != 'e' {
		return nil
	}
	end := i + 2
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - nameOff,
		Replace: "print",
	}}
}

func checkZC1171(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "echo" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}

	first := cmd.Arguments[0].String()
	if first == "-e" {
		return []Violation{{
			KataID: "ZC1171",
			Message: "Use `print` instead of `echo -e`. Zsh `print` natively interprets " +
				"escape sequences and is more portable than `echo -e`.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
