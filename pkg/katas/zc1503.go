// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1503",
		Title:    "Error on `groupadd -g 0` / `groupmod -g 0` — creates duplicate root group",
		Severity: SeverityError,
		Description: "Creating or renaming a group to GID 0 gives its members the same privileges " +
			"as members of `root` for every file that grants permissions to GID 0. Combined " +
			"with `usermod -G 0 <user>` it becomes an invisible privilege escalation path. " +
			"Distro tooling already reserves GID 0 for `root`; pick a sensible unused GID " +
			"(`getent group` gives the list) and scope access via sudoers or polkit.",
		Check: checkZC1503,
	})
}

func checkZC1503(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "groupadd" && ident.Value != "groupmod" && ident.Value != "addgroup" {
		return nil
	}

	var prevG bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevG {
			prevG = false
			if v == "0" {
				return zc1503Violation(cmd)
			}
		}
		if v == "-g" || v == "--gid" {
			prevG = true
		}
		if v == "-g0" || v == "--gid=0" {
			return zc1503Violation(cmd)
		}
	}
	return nil
}

func zc1503Violation(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1503",
		Message: "Creating a group with GID 0 duplicates the `root` group — hidden privesc. " +
			"Pick an unused GID (see `getent group`) and scope via sudoers/polkit.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
