// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1293",
		Title:    "Use `[[ ]]` instead of `test` command in Zsh",
		Severity: SeverityStyle,
		Description: "Zsh `[[ ]]` provides a more powerful conditional expression syntax than " +
			"the `test` command. It supports pattern matching, regex, and does not require " +
			"quoting of variable expansions to prevent word splitting.",
		Check: checkZC1293,
		Fix:   fixZC1293,
	})
}

// fixZC1293 rewrites `test EXPR…` into `[[ EXPR… ]]`. Two edits:
// the `test` command name becomes `[[`; ` ]]` is appended after the
// last argument's source span. Bails when there are no arguments
// (a bare `test` is invalid anyway). Idempotent — a re-run sees
// `[[ ... ]]`, not `test`.
func fixZC1293(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "test" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	cmdOff := LineColToByteOffset(source, v.Line, v.Column)
	if cmdOff < 0 || cmdOff+len("test") > len(source) {
		return nil
	}
	if string(source[cmdOff:cmdOff+len("test")]) != "test" {
		return nil
	}
	lastArg := cmd.Arguments[len(cmd.Arguments)-1]
	lastTok := lastArg.TokenLiteralNode()
	lastOff := LineColToByteOffset(source, lastTok.Line, lastTok.Column)
	if lastOff < 0 {
		return nil
	}
	lastVal := lastArg.String()
	lastEnd := lastOff + len(lastVal)
	if lastEnd > len(source) {
		return nil
	}
	endLine, endCol := offsetLineColZC1293(source, lastEnd)
	if endLine < 0 {
		return nil
	}
	return []FixEdit{
		{Line: v.Line, Column: v.Column, Length: len("test"), Replace: "[["},
		{Line: endLine, Column: endCol, Length: 0, Replace: " ]]"},
	}
}

func offsetLineColZC1293(source []byte, offset int) (int, int) {
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

func checkZC1293(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "test" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1293",
		Message: "Use `[[ ]]` instead of the `test` command in Zsh. `[[ ]]` is more powerful and does not require variable quoting.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityStyle,
	}}
}
