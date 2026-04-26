// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1586",
		Title:    "Style: `chkconfig` / `update-rc.d` / `insserv` — SysV init relics, use `systemctl`",
		Severity: SeverityStyle,
		Description: "`chkconfig` (Red Hat), `update-rc.d` (Debian), and `insserv` (SUSE) are " +
			"SysV-init compatibility wrappers for enabling/disabling services at boot. On any " +
			"distro that has used systemd for the last decade they are translated to " +
			"`systemctl enable|disable`, but silently lose unit-template arguments, " +
			"`[Install]` alias handling, and socket-activated services. Call `systemctl " +
			"enable <unit>` directly.",
		Check: checkZC1586,
	})
}

func checkZC1586(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "chkconfig" && ident.Value != "update-rc.d" &&
		ident.Value != "insserv" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1586",
		Message: "`" + ident.Value + "` is a SysV-init relic. Use `systemctl enable|disable " +
			"<unit>` directly.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
