// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1241",
		Title:    "Use `xargs -0` with null separators for safe argument passing",
		Severity: SeverityWarning,
		Description: "`xargs` without `-0` splits on whitespace, breaking on filenames with spaces. " +
			"Use `xargs -0` paired with `find -print0` for safe handling.",
		Check: checkZC1241,
		Fix:   fixZC1241,
	})
}

// fixZC1241 inserts ` -0` after the `xargs` command name so
// null-terminated input from `find -print0` is consumed safely.
// Detector gates on `xargs rm` without `-0`.
func fixZC1241(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "xargs" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOff)
	if nameLen != len("xargs") {
		return nil
	}
	insertAt := nameOff + nameLen
	insLine, insCol := offsetLineColZC1241(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -0",
	}}
}

func offsetLineColZC1241(source []byte, offset int) (int, int) {
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

func checkZC1241(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "xargs" {
		return nil
	}

	hasNull := false
	hasRM := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-0" {
			hasNull = true
		}
		if val == "rm" {
			hasRM = true
		}
	}

	if hasRM && !hasNull {
		return []Violation{{
			KataID: "ZC1241",
			Message: "Use `xargs -0 rm` with `find -print0` for safe deletion. " +
				"Without `-0`, filenames with spaces or special characters break.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
