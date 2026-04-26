// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1912",
		Title:    "Warn on `dhclient -r` / `dhclient -x` / `dhcpcd -k` — drops the lease and breaks network",
		Severity: SeverityWarning,
		Description: "`dhclient -r` releases the current DHCP lease (sending a DHCPRELEASE), " +
			"`dhclient -x` terminates the daemon without releasing, and `dhcpcd -k` does the " +
			"equivalent for dhcpcd. On a remote host the very next thing that happens is the " +
			"SSH session drops, and in a VPC any automation waiting for a reply never sees " +
			"one. Stage the release together with a re-acquire (`dhclient -1 $iface` or " +
			"`nmcli device reapply $iface`) or schedule it via `systemd-run --on-active=` " +
			"so the operator is not cut off mid-session.",
		Check: checkZC1912,
	})
}

func checkZC1912(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "dhclient":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "-r" || v == "-x" || v == "--release" {
				return zc1912Hit(cmd, "dhclient "+v)
			}
		}
	case "dhcpcd":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "-k" || v == "--release" {
				return zc1912Hit(cmd, "dhcpcd "+v)
			}
		}
	}
	return nil
}

func zc1912Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1912",
		Message: "`" + form + "` drops the DHCP lease — SSH session cuts, VPC " +
			"reachability stalls. Pair with a re-acquire (`dhclient -1`/`nmcli device " +
			"reapply`), or schedule via `systemd-run --on-active=`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
