// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1210",
		Title:    "Use `journalctl --no-pager` in scripts",
		Severity: SeverityStyle,
		Description: "`journalctl` invokes a pager by default which hangs in non-interactive scripts. " +
			"Use `--no-pager` for reliable script output.",
		Check: checkZC1210,
		Fix:   fixZC1210,
	})
}

// fixZC1210 inserts ` --no-pager` after the `journalctl` command
// name, preventing the pager from hanging in non-interactive runs.
// Mirrors ZC1209's insertion for `systemctl`.
func fixZC1210(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "journalctl" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOff)
	if nameLen != len("journalctl") {
		return nil
	}
	insertAt := nameOff + nameLen
	insLine, insCol := offsetLineColZC1210(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " --no-pager",
	}}
}

func offsetLineColZC1210(source []byte, offset int) (int, int) {
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

func checkZC1210(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "journalctl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "--no-pager" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1210",
		Message: "Use `journalctl --no-pager` in scripts. Without it, " +
			"journalctl invokes a pager that hangs in non-interactive execution.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
