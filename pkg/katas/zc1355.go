// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1355",
		Title:    "Use `print -r` instead of `echo -E` for raw output",
		Severity: SeverityStyle,
		Description: "`echo -E` disables backslash interpretation, but the flag is Bash-ism and " +
			"ignored by POSIX `echo`. Zsh's `print -r` is the idiomatic raw-printer; combine " +
			"with `-n` (no newline), `-l` (one per line), `-u<fd>` (file descriptor), or `--` " +
			"(end of flags) as needed.",
		Check: checkZC1355,
		Fix:   fixZC1355,
	})
}

// fixZC1355 collapses `echo -E` into `print -r`. Span covers the
// command name, intervening whitespace, and the `-E` flag.
func fixZC1355(node ast.Node, v Violation, source []byte) []FixEdit {
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
	if i+2 > len(source) || source[i] != '-' || source[i+1] != 'E' {
		return nil
	}
	end := i + 2
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - nameOff,
		Replace: "print -r",
	}}
}

func checkZC1355(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "echo" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-E" {
			return []Violation{{
				KataID: "ZC1355",
				Message: "Use `print -r` instead of `echo -E` for raw output. " +
					"`-E` is a Bash-ism and ignored by POSIX echo.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
