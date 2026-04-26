// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1170",
		Title:    "Avoid `pushd`/`popd` without `-q` flag",
		Severity: SeverityStyle,
		Description: "`pushd` and `popd` print the directory stack by default, cluttering output. " +
			"Use `-q` flag to suppress output in scripts.",
		Check: checkZC1170,
		Fix:   fixZC1170,
	})
}

// fixZC1170 inserts ` -q` after `pushd` or `popd` so the directory
// stack output is suppressed in scripts.
func fixZC1170(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || (ident.Value != "pushd" && ident.Value != "popd") {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOff)
	if nameLen != len(ident.Value) {
		return nil
	}
	insertAt := nameOff + nameLen
	insLine, insCol := offsetLineColZC1170(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -q",
	}}
}

func offsetLineColZC1170(source []byte, offset int) (int, int) {
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

func checkZC1170(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "pushd" && ident.Value != "popd" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-q" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1170",
		Message: "Use `" + ident.Value + " -q` to suppress directory stack output in scripts. " +
			"Without `-q`, the stack is printed on every call.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
