// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1798",
		Title:    "Warn on `ufw reset` — wipes every configured firewall rule",
		Severity: SeverityWarning,
		Description: "`ufw reset` returns the firewall to the distro default: every user-defined " +
			"rule is removed, default incoming policy reverts (usually to `deny`, but the net " +
			"effect is the loss of every allow-list entry the host relied on). Paired with " +
			"`--force`, no prompt is issued. In a provisioning script the operation is " +
			"sometimes desired to start from a clean slate, but running it mid-session or on " +
			"a host that currently serves traffic drops connections without warning. Snapshot " +
			"the rules first (`ufw status numbered > /tmp/ufw.bak`), and prefer removing " +
			"specific rules with `ufw delete <num>` over a full reset.",
		Check: checkZC1798,
	})
}

func checkZC1798(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `ufw --force reset` mangles to name=`force` with `reset` as arg[0].
	if ident.Value == "force" {
		if len(cmd.Arguments) > 0 && cmd.Arguments[0].String() == "reset" {
			return zc1798Hit(cmd, "ufw --force reset")
		}
		return nil
	}

	if ident.Value != "ufw" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	if cmd.Arguments[0].String() != "reset" {
		return nil
	}
	return zc1798Hit(cmd, "ufw reset")
}

func zc1798Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1798",
		Message: "`" + what + "` drops every user-defined firewall rule. Snapshot " +
			"(`ufw status numbered > /tmp/ufw.bak`) first, and prefer " +
			"`ufw delete <num>` for targeted removals.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
