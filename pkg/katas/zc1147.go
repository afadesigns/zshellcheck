// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1147",
		Title:    "Avoid `mkdir` without `-p` for nested paths",
		Severity: SeverityInfo,
		Description: "Using `mkdir` without `-p` fails if parent directories don't exist. " +
			"Use `mkdir -p` to create the full path safely.",
		Check: checkZC1147,
		Fix:   fixZC1147,
	})
}

// fixZC1147 inserts ` -p` after the `mkdir` command name so nested
// paths survive missing intermediates. Detector already gates on
// absence of `-p` and presence of a nested path.
func fixZC1147(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "mkdir" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOff)
	if nameLen != len("mkdir") {
		return nil
	}
	insertAt := nameOff + nameLen
	insLine, insCol := offsetLineColZC1147(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -p",
	}}
}

func offsetLineColZC1147(source []byte, offset int) (int, int) {
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

func checkZC1147(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "mkdir" {
		return nil
	}

	hasParentFlag := false
	hasNestedPath := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-p" {
			hasParentFlag = true
		}
		// Check for paths with multiple slashes (nested)
		if len(val) > 0 && val[0] != '-' {
			slashCount := 0
			for _, ch := range val {
				if ch == '/' {
					slashCount++
				}
			}
			if slashCount >= 2 {
				hasNestedPath = true
			}
		}
	}

	if hasParentFlag || !hasNestedPath {
		return nil
	}

	return []Violation{{
		KataID: "ZC1147",
		Message: "Use `mkdir -p` when creating nested directories. " +
			"Without `-p`, `mkdir` fails if parent directories don't exist.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}
