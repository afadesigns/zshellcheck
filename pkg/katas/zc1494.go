// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1494",
		Title:    "Warn on `tcpdump -w <file>` without `-Z <user>` — capture file owned by root",
		Severity: SeverityWarning,
		Description: "`tcpdump` needs root (or CAP_NET_RAW) to open the raw socket, but once the " +
			"socket is open it should drop privileges with `-Z <user>` before writing the pcap. " +
			"Without `-Z`, the capture file is owned by root, any bpf filter bug is exercised " +
			"with root privileges, and on a shared host the pcap can land with permissions that " +
			"leak sensitive traffic to other users. Pair `-w` with `-Z tcpdump` (or a dedicated " +
			"capture user).",
		Check: checkZC1494,
	})
}

func checkZC1494(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "tcpdump" {
		return nil
	}

	hasW := false
	hasZ := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-w" || v == "--write-file" {
			hasW = true
		}
		if v == "-Z" || v == "--relinquish-privileges" {
			hasZ = true
		}
	}
	if !hasW || hasZ {
		return nil
	}
	return []Violation{{
		KataID: "ZC1494",
		Message: "`tcpdump -w` without `-Z <user>` writes the pcap as root and never drops " +
			"privileges. Add `-Z tcpdump` (or a dedicated capture user).",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
