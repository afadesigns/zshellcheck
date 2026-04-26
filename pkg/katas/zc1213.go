// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1213",
		Title:    "Use `apt-get -y` in scripts for non-interactive installs",
		Severity: SeverityWarning,
		Description: "`apt-get install` without `-y` prompts for confirmation which hangs scripts. " +
			"Use `-y` or set `DEBIAN_FRONTEND=noninteractive` for unattended installs.",
		Check: checkZC1213,
		Fix:   fixZC1213,
	})
}

// fixZC1213 inserts ` -y` after `apt-get` so install / upgrade /
// dist-upgrade run without interactive confirmation. Detector
// already guards the shape (install-class subcommand + no -y).
func fixZC1213(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "apt-get" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOff)
	if nameLen != len("apt-get") {
		return nil
	}
	insertAt := nameOff + nameLen
	insLine, insCol := offsetLineColZC1213(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -y",
	}}
}

func offsetLineColZC1213(source []byte, offset int) (int, int) {
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

func checkZC1213(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "apt-get" {
		return nil
	}

	hasInstall := false
	hasYes := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "install" || val == "upgrade" || val == "dist-upgrade" {
			hasInstall = true
		}
		if val == "-y" || val == "--yes" || val == "-qq" {
			hasYes = true
		}
	}

	if hasInstall && !hasYes {
		return []Violation{{
			KataID: "ZC1213",
			Message: "Use `apt-get -y` in scripts. Without `-y`, apt-get prompts for confirmation " +
				"which hangs in non-interactive execution.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
