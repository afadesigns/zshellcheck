package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1601",
		Title:    "Warn on `ethtool -s $IF wol <g|u|m|b|a>` — enables remote Wake-on-LAN",
		Severity: SeverityWarning,
		Description: "Wake-on-LAN powers the host on from a sleep / soft-off state when a " +
			"matching packet reaches the NIC. The wake logic fires in a privileged firmware " +
			"path long before the kernel boots and firewall rules are loaded — so any packet " +
			"that reaches the interface (magic-packet, unicast, broadcast, ARP) triggers the " +
			"power-on unfiltered. On a shared or public LAN attackers on the broadcast domain " +
			"can wake hosts at will. Keep `wol d` (disable) unless a documented operational " +
			"need requires one of the wake bits.",
		Check: checkZC1601,
	})
}

func checkZC1601(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ethtool" {
		return nil
	}
	if len(cmd.Arguments) < 4 {
		return nil
	}
	if cmd.Arguments[0].String() != "-s" {
		return nil
	}

	for i := 2; i+1 < len(cmd.Arguments); i++ {
		if cmd.Arguments[i].String() != "wol" {
			continue
		}
		bits := cmd.Arguments[i+1].String()
		if bits == "d" {
			return nil
		}
		enables := false
		for _, c := range bits {
			switch c {
			case 'g', 'u', 'm', 'b', 'a', 'p', 's', 'f':
				enables = true
			}
		}
		if !enables {
			return nil
		}
		return []Violation{{
			KataID: "ZC1601",
			Message: "`ethtool -s " + cmd.Arguments[1].String() + " wol " + bits + "` " +
				"enables Wake-on-LAN — the NIC powers the host on before firewall rules " +
				"load. Keep `wol d` unless a documented operational need requires " +
				strings.TrimSpace(bits) + ".",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}
