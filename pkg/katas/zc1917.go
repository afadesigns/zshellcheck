// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1917",
		Title:    "Info on `iw dev $IF scan` / `iwlist $IF scan` — active WiFi scan from a script",
		Severity: SeverityInfo,
		Description: "`iw dev wlan0 scan` (and the older `iwlist wlan0 scan`) performs an active " +
			"probe-request sweep across every supported channel. It requires `CAP_NET_ADMIN`, " +
			"briefly interrupts the current association, and announces the host's presence to " +
			"every nearby access point — logs on the other side will show one MAC asking " +
			"about every SSID. Use the cached `iw dev $IF link` / `iwctl station $IF show` " +
			"for passive lookups, and reserve `scan` for diagnostic sessions with console " +
			"approval rather than background scripts.",
		Check: checkZC1917,
	})
}

func checkZC1917(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "iw":
		if zc1917HasSubcmd(cmd, "scan") {
			return zc1917Hit(cmd, "iw dev <if> scan")
		}
	case "iwlist":
		if zc1917HasSubcmd(cmd, "scan") || zc1917HasSubcmd(cmd, "scanning") {
			return zc1917Hit(cmd, "iwlist <if> scan")
		}
	}
	return nil
}

func zc1917HasSubcmd(cmd *ast.SimpleCommand, sub string) bool {
	for _, arg := range cmd.Arguments {
		if arg.String() == sub {
			return true
		}
	}
	return false
}

func zc1917Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1917",
		Message: "`" + form + "` runs an active probe-request sweep — interrupts the " +
			"current association and broadcasts the host to every nearby AP. Use cached " +
			"`iw dev $IF link` for passive queries.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}
