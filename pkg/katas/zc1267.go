// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1267",
		Title:    "Use `df -P` for POSIX-portable disk usage output",
		Severity: SeverityStyle,
		Description: "`df -h` output format varies across systems and locales. " +
			"Use `df -P` for single-line, fixed-format output safe for script parsing.",
		Check: checkZC1267,
		Fix:   fixZC1267,
	})
}

// fixZC1267 inserts ` -P` after the `df` command name. Detector
// narrows to `df -h` (script-unsafe), so only that shape is
// rewritten.
func fixZC1267(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "df" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOff)
	if nameLen != len("df") {
		return nil
	}
	insertAt := nameOff + nameLen
	insLine, insCol := offsetLineColZC1267(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -P",
	}}
}

func offsetLineColZC1267(source []byte, offset int) (int, int) {
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

func checkZC1267(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "df" {
		return nil
	}

	hasPortable := false
	hasHuman := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-P" {
			hasPortable = true
		}
		if val == "-h" {
			hasHuman = true
		}
	}

	if hasHuman && !hasPortable {
		return []Violation{{
			KataID: "ZC1267",
			Message: "Use `df -P` for script-safe output. `df -h` format varies across " +
				"systems and may split long device names across lines.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
