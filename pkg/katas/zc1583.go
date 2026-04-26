// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1583",
		Title:    "Warn on `find ... -delete` without `-maxdepth` — unbounded recursive delete",
		Severity: SeverityWarning,
		Description: "`find PATH -delete` walks the tree recursively and removes every match. " +
			"Without `-maxdepth N` the walk crosses into every subtree, including symlinks " +
			"that point outside the intended scope and mount points that expand the blast " +
			"radius. Scope the depth (`-maxdepth 2`) and prefer a dry-run first " +
			"(`find ... -print | head`).",
		Check: checkZC1583,
	})
}

func checkZC1583(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "find" {
		return nil
	}

	var hasDelete, hasMaxdepth, hasPrune, hasXdev bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "-delete":
			hasDelete = true
		case "-maxdepth":
			hasMaxdepth = true
		case "-prune":
			hasPrune = true
		case "-xdev", "-mount":
			hasXdev = true
		}
	}
	if !hasDelete || hasMaxdepth || hasPrune || hasXdev {
		return nil
	}
	return []Violation{{
		KataID: "ZC1583",
		Message: "`find -delete` without `-maxdepth` / `-xdev` / `-prune` walks the whole " +
			"tree. Scope the depth (e.g. `-maxdepth 2`) and dry-run first.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
