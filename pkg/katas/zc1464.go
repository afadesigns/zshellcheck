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

func checkZC1464(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "iptables" && ident.Value != "ip6tables" && ident.Value != "nft" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, arg := range cmd.Arguments {
		args = append(args, arg.String())
	}

	// Catch: `iptables -F` (no chain) — flushes everything
	// Catch: `iptables -P INPUT ACCEPT` / `-P FORWARD ACCEPT`
	for i, a := range args {
		if a == "-F" || a == "--flush" {
			return violateZC1464(cmd, "flushing all firewall rules")
		}
		if (a == "-P" || a == "--policy") && i+2 < len(args) && args[i+2] == "ACCEPT" {
			chain := args[i+1]
			if chain == "INPUT" || chain == "FORWARD" {
				return violateZC1464(cmd, "default-ACCEPT policy on "+chain+" chain")
			}
		}
	}
	return nil
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
