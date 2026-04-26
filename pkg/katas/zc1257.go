// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1257",
		Title:    "Use `docker stop -t` to set graceful shutdown timeout",
		Severity: SeverityStyle,
		Description: "`docker stop` defaults to 10s before SIGKILL. In CI scripts, " +
			"set an explicit timeout with `-t` to control shutdown behavior.",
		Check: checkZC1257,
		Fix:   fixZC1257,
	})
}

// fixZC1257 inserts ` -t 10` after the `stop` subcommand of a
// `docker stop …` invocation. Mirrors the subcommand-level pattern
// used by ZC1265 (`systemctl enable --now`).
func fixZC1257(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}
	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "stop" {
		return nil
	}
	stopArg := cmd.Arguments[0]
	tok := stopArg.TokenLiteralNode()
	off := LineColToByteOffset(source, tok.Line, tok.Column)
	if off < 0 || off+len("stop") > len(source) {
		return nil
	}
	if string(source[off:off+len("stop")]) != "stop" {
		return nil
	}
	insertAt := off + len("stop")
	insLine, insCol := offsetLineColZC1257(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -t 10",
	}}
}

func offsetLineColZC1257(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1257(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "stop" {
		return nil
	}

	hasTimeout := false
	for _, arg := range cmd.Arguments[1:] {
		if arg.String() == "-t" {
			hasTimeout = true
		}
	}

	if !hasTimeout {
		return []Violation{{
			KataID: "ZC1257",
			Message: "Use `docker stop -t N` to set an explicit shutdown timeout. " +
				"The default 10s may be too long or too short for your use case.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
