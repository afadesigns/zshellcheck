package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1466",
		Title:    "Warn on disabling the host firewall (`ufw disable` / `systemctl stop firewalld`)",
		Severity: SeverityWarning,
		Description: "Disabling the host firewall leaves every listening port reachable from " +
			"every network the host is on. This is a common \"just make it work\" shortcut that " +
			"has shipped to production more than once. Keep the firewall running and open the " +
			"specific port with `ufw allow <port>` / `firewall-cmd --add-port=<port>/tcp`.",
		Check: checkZC1466,
	})
}

func checkZC1466(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value == "ufw" && len(cmd.Arguments) >= 1 {
		if cmd.Arguments[0].String() == "disable" {
			return violateZC1466(cmd, "ufw disable")
		}
	}

	if ident.Value == "systemctl" && len(cmd.Arguments) >= 2 {
		verb := cmd.Arguments[0].String()
		unit := cmd.Arguments[1].String()
		if (verb == "stop" || verb == "disable" || verb == "mask") &&
			(unit == "firewalld" || unit == "firewalld.service" ||
				unit == "ufw" || unit == "ufw.service" ||
				unit == "nftables" || unit == "nftables.service" ||
				unit == "iptables" || unit == "iptables.service") {
			return violateZC1466(cmd, "systemctl "+verb+" "+unit)
		}
	}

	return nil
}

func violateZC1466(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID:  "ZC1466",
		Message: "Host firewall disabled (" + what + "). Keep it on and open specific ports.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityWarning,
	}}
}
