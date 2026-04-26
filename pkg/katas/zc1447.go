// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1447",
		Title:    "Avoid deprecated `ifconfig` / `netstat` — prefer `ip` / `ss`",
		Severity: SeverityStyle,
		Description: "On modern Linux, `ifconfig` and `netstat` (from net-tools) are deprecated " +
			"in favor of the iproute2 suite: `ip addr`, `ip link`, `ip route`, `ss`. net-tools " +
			"is not installed by default on many distros (Alpine, Fedora Cloud, minimal images), " +
			"so scripts break. Use iproute2 commands for portability.",
		Check: checkZC1447,
	})
}

func checkZC1447(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "ifconfig":
		return []Violation{{
			KataID:  "ZC1447",
			Message: "`ifconfig` is deprecated. Use `ip addr` / `ip link` / `ip route` from iproute2.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	case "netstat":
		return []Violation{{
			KataID:  "ZC1447",
			Message: "`netstat` is deprecated. Use `ss` from iproute2 (same flags, faster output).",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	case "route":
		return []Violation{{
			KataID:  "ZC1447",
			Message: "`route` is deprecated. Use `ip route`.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}
