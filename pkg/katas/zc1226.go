// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1226",
		Title:    "Use `dmesg -T` or `--time-format=iso` for readable timestamps",
		Severity: SeverityStyle,
		Description: "`dmesg` without `-T` shows raw kernel timestamps in seconds since boot. " +
			"Use `-T` for human-readable timestamps or `--time-format=iso` for ISO 8601.",
		Check: checkZC1226,
		Fix:   fixZC1226,
	})
}

// fixZC1226 inserts ` -T` after the `dmesg` command name. Mirrors
// other insertion-style fixes (ZC1012 / ZC1017 / ZC1170 / ZC1209).
func fixZC1226(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "dmesg" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOff)
	if nameLen != len("dmesg") {
		return nil
	}
	insertAt := nameOff + nameLen
	insLine, insCol := offsetLineColZC1226(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -T",
	}}
}

func offsetLineColZC1226(source []byte, offset int) (int, int) {
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

func checkZC1226(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "dmesg" {
		return nil
	}

	hasTimeFlag := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-T" || val == "-t" || val == "--ctime" || val == "--reltime" {
			hasTimeFlag = true
		}
	}

	if !hasTimeFlag && len(cmd.Arguments) > 0 {
		return []Violation{{
			KataID: "ZC1226",
			Message: "Use `dmesg -T` for human-readable timestamps instead of raw " +
				"kernel boot-seconds. Or use `--time-format=iso` for ISO 8601.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
