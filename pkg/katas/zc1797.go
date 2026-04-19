package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1797",
		Title:    "Warn on `ip link set <iface> down` / `ifdown <iface>` — locks out remote admin on that path",
		Severity: SeverityWarning,
		Description: "Taking a network interface down from an SSH session that rides on the same " +
			"interface cuts the script off mid-run: the TCP connection freezes, any later " +
			"step silently fails, and recovery requires console / out-of-band access. " +
			"Common bugs are typos (`eth1` instead of `eth0`), scripts that target the only " +
			"uplink on a cloud VM, or running the command without first confirming that the " +
			"interface is not the one carrying the admin session. Wrap the `down` in a " +
			"`systemd-run --on-active=30s --unit=recover ip link set <iface> up` rollback, " +
			"or stage both `down` and `up` through `nmcli connection up/down` with a pinned " +
			"fallback profile.",
		Check: checkZC1797,
	})
}

func checkZC1797(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "ip":
		// Expect: ip link set IFACE down
		if len(cmd.Arguments) < 4 {
			return nil
		}
		if cmd.Arguments[0].String() != "link" || cmd.Arguments[1].String() != "set" {
			return nil
		}
		for _, arg := range cmd.Arguments[2:] {
			if arg.String() == "down" {
				return zc1797Hit(cmd, "ip link set … down")
			}
		}
	case "ifdown":
		if len(cmd.Arguments) == 0 {
			return nil
		}
		first := cmd.Arguments[0].String()
		if first == "--help" || first == "-h" {
			return nil
		}
		return zc1797Hit(cmd, "ifdown "+first)
	}
	return nil
}

func zc1797Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1797",
		Message: "`" + what + "` disables a network interface — if it carries the " +
			"SSH session, the script cuts itself off. Schedule a rollback via " +
			"`systemd-run --on-active=30s ip link set … up` or stage via `nmcli`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
