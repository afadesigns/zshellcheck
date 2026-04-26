// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1464",
		Title:    "Warn on `iptables -F` / `-P INPUT ACCEPT` — flushes or opens the host firewall",
		Severity: SeverityWarning,
		Description: "Flushing all rules (`-F`) or setting the default INPUT/FORWARD policy to " +
			"ACCEPT leaves the host with no network filter. This is rarely correct outside a " +
			"first-boot provisioning script, and is a frequent post-compromise persistence step. " +
			"Use `iptables-save`/`iptables-restore` for atomic reloads and keep a default-drop " +
			"policy on all hook chains.",
		Check: checkZC1464,
	})
}

var (
	zc1464FirewallNames = map[string]struct{}{"iptables": {}, "ip6tables": {}, "nft": {}}
	zc1464FlushFlags    = map[string]struct{}{"-F": {}, "--flush": {}}
	zc1464PolicyFlags   = map[string]struct{}{"-P": {}, "--policy": {}}
	zc1464OpenChains    = map[string]struct{}{"INPUT": {}, "FORWARD": {}}
)

func checkZC1464(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if _, hit := zc1464FirewallNames[CommandIdentifier(cmd)]; !hit {
		return nil
	}
	args := zc1464StringArgs(cmd)
	if zc1464FlushHit(args) {
		return violateZC1464(cmd, "flushing all firewall rules")
	}
	if chain := zc1464AcceptPolicy(args); chain != "" {
		return violateZC1464(cmd, "default-ACCEPT policy on "+chain+" chain")
	}
	return nil
}

func zc1464StringArgs(cmd *ast.SimpleCommand) []string {
	out := make([]string, 0, len(cmd.Arguments))
	for _, arg := range cmd.Arguments {
		out = append(out, arg.String())
	}
	return out
}

func zc1464FlushHit(args []string) bool {
	for _, a := range args {
		if _, hit := zc1464FlushFlags[a]; hit {
			return true
		}
	}
	return false
}

func zc1464AcceptPolicy(args []string) string {
	for i, a := range args {
		if _, hit := zc1464PolicyFlags[a]; !hit {
			continue
		}
		if i+2 >= len(args) || args[i+2] != "ACCEPT" {
			continue
		}
		if _, hit := zc1464OpenChains[args[i+1]]; hit {
			return args[i+1]
		}
	}
	return ""
}

func violateZC1464(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID:  "ZC1464",
		Message: "Firewall hardening weakened (" + what + "). Keep default-drop and use atomic reload.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityWarning,
	}}
}
