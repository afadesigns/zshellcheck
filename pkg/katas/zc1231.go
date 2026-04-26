// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1231",
		Title:    "Use `git clone --depth 1` for CI and build scripts",
		Severity: SeverityStyle,
		Description: "`git clone` without `--depth` downloads the entire history. " +
			"Use `--depth 1` in CI/build scripts where only the latest commit is needed.",
		Check: checkZC1231,
		Fix:   fixZC1231,
	})
}

// fixZC1231 inserts ` --depth 1` after the `clone` subcommand in
// `git clone …`. Mirrors ZC1234's subcommand-level insertion for
// docker run --rm.
func fixZC1231(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	cloneArg := cmd.Arguments[0]
	if cloneArg.String() != "clone" {
		return nil
	}
	tok := cloneArg.TokenLiteralNode()
	off := LineColToByteOffset(source, tok.Line, tok.Column)
	if off < 0 || off+5 > len(source) {
		return nil
	}
	if string(source[off:off+5]) != "clone" {
		return nil
	}
	insertAt := off + 5
	insLine, insCol := offsetLineColZC1231(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " --depth 1",
	}}
}

func offsetLineColZC1231(source []byte, offset int) (int, int) {
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

func checkZC1231(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}

	if cmd.Arguments[0].String() != "clone" {
		return nil
	}

	hasDepth := false
	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "--depth" || val == "--shallow-since" || val == "--single-branch" {
			hasDepth = true
		}
	}

	if !hasDepth {
		return []Violation{{
			KataID: "ZC1231",
			Message: "Consider `git clone --depth 1` in scripts. Full clones download " +
				"entire history which is unnecessary for builds and CI.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
