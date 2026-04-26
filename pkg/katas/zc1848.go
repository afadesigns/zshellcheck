// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1848",
		Title:    "Warn on `ssh -o CheckHostIP=no` — DNS-spoof warning for known hosts silenced",
		Severity: SeverityWarning,
		Description: "`CheckHostIP` (on by default) stores the host's IP address alongside its " +
			"host key in `~/.ssh/known_hosts`; if DNS later resolves the same name to a " +
			"different IP but the key still matches, ssh warns you. Turning the check off " +
			"with `-o CheckHostIP=no` keeps the host-key comparison but silences the " +
			"IP-mismatch warning — which means a DNS-poisoning attacker who already holds " +
			"the previously-seen host key (stolen, misplaced backup, leaked by a " +
			"decommissioned box) can route the session through their box without a peep. " +
			"Leave the default, and if you really need to skip the IP record (load-balanced " +
			"pool with shared keys) document the risk and prefer `HostKeyAlias` instead.",
		Check: checkZC1848,
	})
}

func checkZC1848(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ssh" && ident.Value != "scp" && ident.Value != "sftp" {
		return nil
	}
	args := cmd.Arguments
	for i, arg := range args {
		v := arg.String()
		var kv string
		switch {
		case v == "-o" && i+1 < len(args):
			kv = args[i+1].String()
		case strings.HasPrefix(v, "-o"):
			kv = strings.TrimPrefix(v, "-o")
		default:
			continue
		}
		if zc1848IsCheckHostIPNo(kv) {
			return []Violation{{
				KataID: "ZC1848",
				Message: "`" + ident.Value + " -o CheckHostIP=no` silences the " +
					"IP-mismatch warning for known hosts — a DNS-spoof + leaked " +
					"host-key attack goes undetected. Leave the default, or use " +
					"`HostKeyAlias` for load-balanced pools.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func zc1848IsCheckHostIPNo(kv string) bool {
	norm := strings.ToLower(strings.Trim(kv, "\"' \t"))
	for _, frag := range []string{"checkhostip=no", "checkhostip = no", "checkhostip=false", "checkhostip=0", "checkhostip=off"} {
		if norm == frag {
			return true
		}
	}
	// Tolerate stray spaces around `=`.
	if strings.HasPrefix(norm, "checkhostip") {
		rest := strings.TrimPrefix(norm, "checkhostip")
		rest = strings.TrimSpace(rest)
		if strings.HasPrefix(rest, "=") {
			val := strings.TrimSpace(strings.TrimPrefix(rest, "="))
			return val == "no" || val == "false" || val == "0" || val == "off"
		}
	}
	return false
}
