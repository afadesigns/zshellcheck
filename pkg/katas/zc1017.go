// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1017",
		Title: "Use `print -r` to print strings literally",
		Description: "The `print` command interprets backslash escape sequences by default. " +
			"To print a string literally, use the `-r` option.",
		Severity: SeverityStyle,
		Check:    checkZC1017,
		Fix:      fixZC1017,
	})
}

// fixZC1017 inserts ` -r` directly after the `print` command name.
// Existing flags are left in place, mirroring ZC1012's `read -r`
// insertion: `print "x"` becomes `print -r "x"`, `print -n "x"`
// becomes `print -r -n "x"`. Idempotent on a second pass — once
// `-r` appears among the flags the detector no longer fires.
func fixZC1017(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	name, ok := cmd.Name.(*ast.Identifier)
	if !ok || name.Value != "print" {
		return nil
	}
	nameOffset := LineColToByteOffset(source, v.Line, v.Column)
	if nameOffset < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOffset)
	if nameLen != len("print") {
		return nil
	}
	insertAt := nameOffset + nameLen
	insLine, insCol := byteOffsetToLineColZC1017(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -r",
	}}
}

func byteOffsetToLineColZC1017(source []byte, offset int) (int, int) {
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

func checkZC1017(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "print" {
			hasRFlag := false
			for _, arg := range cmd.Arguments {
				argStr := arg.String()
				argStr = strings.Trim(argStr, "\"'")
				if strings.HasPrefix(argStr, "-") && strings.Contains(argStr, "r") {
					hasRFlag = true
					break
				}
			}
			if !hasRFlag {
				violations = append(violations, Violation{
					KataID:  "ZC1017",
					Message: "Use `print -r` to print strings literally.",
					Line:    name.Token.Line,
					Column:  name.Token.Column,
					Level:   SeverityStyle,
				})
			}
		}
	}

	return violations
}
