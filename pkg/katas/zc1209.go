// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1209",
		Title:    "Use `systemctl --no-pager` in scripts",
		Severity: SeverityStyle,
		Description: "`systemctl` invokes a pager by default which hangs in non-interactive scripts. " +
			"Use `--no-pager` or pipe to `cat` for reliable script output.",
		Check: checkZC1209,
		Fix:   fixZC1209,
	})
}

// fixZC1209 inserts ` --no-pager` after the `systemctl` command
// name so subcommands that emit pager output (status, list-*)
// behave predictably in scripts.
func fixZC1209(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "systemctl" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOff)
	if nameLen != len("systemctl") {
		return nil
	}
	insertAt := nameOff + nameLen
	insLine, insCol := offsetLineColZC1209(source, insertAt)
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

func offsetLineColZC1209(source []byte, offset int) (int, int) {
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

func checkZC1209(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "systemctl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "--no-pager" {
			return nil
		}
	}

	// Only flag subcommands that produce output (status, list-units, etc.)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "status" || val == "list-units" || val == "list-timers" || val == "show" {
			return []Violation{{
				KataID: "ZC1209",
				Message: "Use `systemctl --no-pager` in scripts. Without it, " +
					"systemctl invokes a pager that hangs in non-interactive execution.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}
