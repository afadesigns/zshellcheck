// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1118",
		Title: "Use `print -rn` instead of `echo -n`",
		Description: "The behavior of `echo -n` varies across shells and platforms. " +
			"In Zsh, `print -rn` is the reliable way to output text without a trailing newline.",
		Severity: SeverityStyle,
		Check:    checkZC1118,
		Fix:      fixZC1118,
	})
}

// fixZC1118 collapses `echo -n` (with any whitespace between) into
// `print -rn`. Spans the `echo` name, intervening whitespace, and
// the `-n` flag in a single edit; remaining arguments stay in place.
func fixZC1118(node ast.Node, v Violation, source []byte) []FixEdit {
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
	// Walk forward past whitespace, then expect `-n`.
	i := nameOff + len("echo")
	for i < len(source) && (source[i] == ' ' || source[i] == '\t') {
		i++
	}
	if i+2 > len(source) || source[i] != '-' || source[i+1] != 'n' {
		return nil
	}
	end := i + 2
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - nameOff,
		Replace: "print -rn",
	}}
}

func checkZC1118(node ast.Node) []Violation {
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

	firstArg := cmd.Arguments[0].String()
	if firstArg == "-n" {
		return []Violation{{
			KataID: "ZC1118",
			Message: "Use `print -rn` instead of `echo -n`. " +
				"`echo -n` behavior varies across shells; `print -rn` is the reliable Zsh idiom.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
