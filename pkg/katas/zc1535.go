package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1535",
		Title:    "Warn on `ip link set <iface> promisc on` — enables packet capture",
		Severity: SeverityWarning,
		Description: "Putting an interface into promiscuous mode tells the NIC to deliver every " +
			"frame to userspace, not just frames addressed to this host. Legitimate for tools " +
			"like tcpdump/tshark (which turn it on themselves) but running it from a script " +
			"and leaving it on is a sniffer-in-place — traffic from other hosts on the same " +
			"broadcast domain lands in anyone's `tshark -i`. Re-disable as soon as capture is " +
			"done, and prefer giving tcpdump `CAP_NET_RAW` so the mode is scoped to a single " +
			"invocation.",
		Check: checkZC1535,
	})
}

func checkZC1535(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ip" && ident.Value != "ifconfig" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	// ip link set <iface> promisc on
	if ident.Value == "ip" {
		for i := 0; i+4 < len(args); i++ {
			if args[i] == "link" && args[i+1] == "set" && args[i+3] == "promisc" && args[i+4] == "on" {
				return zc1535Violation(cmd)
			}
		}
	}
	// ifconfig <iface> promisc
	if ident.Value == "ifconfig" {
		for _, a := range args {
			if a == "promisc" {
				return zc1535Violation(cmd)
			}
		}
	}
	return nil
}

func zc1535Violation(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1535",
		Message: "Interface put into promiscuous mode — sniffer-in-place. Re-disable after " +
			"capture, or grant tcpdump CAP_NET_RAW instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
