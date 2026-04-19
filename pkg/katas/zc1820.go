package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1820",
		Title:    "Warn on `netplan apply` — applies network config immediately with no rollback timer",
		Severity: SeverityWarning,
		Description: "`netplan apply` regenerates the rendered backend config (systemd-networkd " +
			"or NetworkManager) and brings it live right away. A mistake in the YAML — wrong " +
			"interface name, missing `dhcp4`, bad addresses, conflicting routes — drops the " +
			"admin SSH session, and recovery needs console access. Run `netplan try` first: " +
			"it applies the new config, waits for confirmation, and rolls back automatically " +
			"if no keypress arrives within the timeout. Only fall through to `netplan apply` " +
			"after the try window has elapsed successfully.",
		Check: checkZC1820,
	})
}

func checkZC1820(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "netplan" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	if cmd.Arguments[0].String() != "apply" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1820",
		Message: "`netplan apply` commits the YAML immediately — a mistake drops the " +
			"admin SSH session with no automatic rollback. Run `netplan try` first " +
			"(auto-reverts if no keypress within the timeout), then `netplan apply`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
