// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1506",
		Title:    "Warn on `newgrp <group>` in scripts — spawns a new shell, breaks control flow",
		Severity: SeverityWarning,
		Description: "`newgrp` starts a new login shell with the requested primary group. Inside " +
			"a non-interactive script that shell inherits no commands, so the script either " +
			"hangs waiting for the user or exits immediately depending on stdin. If the script " +
			"genuinely needs temporarily-augmented group access, call `sg <group> -c <cmd>` " +
			"or, in a service context, use `SupplementaryGroups=` in the unit file.",
		Check: checkZC1506,
	})
}

func checkZC1506(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "newgrp" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1506",
		Message: "`newgrp` starts a new shell — script either hangs or exits. Use " +
			"`sg <group> -c <cmd>` or systemd `SupplementaryGroups=`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
