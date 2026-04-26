// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1242",
		Title:    "Use `tar -C dir` to extract into a specific directory",
		Severity: SeverityInfo,
		Description: "`tar xf` without `-C` extracts into the current directory which may " +
			"overwrite files unexpectedly. Use `-C dir` to control the extraction target.",
		Check: checkZC1242,
	})
}

func checkZC1242(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "tar" {
		return nil
	}

	hasExtract := false
	hasTarget := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-x" || val == "xf" || val == "xzf" || val == "xjf" || val == "xJf" {
			hasExtract = true
		}
		if val == "-C" {
			hasTarget = true
		}
	}

	if hasExtract && !hasTarget {
		return []Violation{{
			KataID: "ZC1242",
			Message: "Use `tar -C dir` to specify extraction directory. " +
				"Without `-C`, tar extracts into the current directory which may overwrite files.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityInfo,
		}}
	}

	return nil
}
