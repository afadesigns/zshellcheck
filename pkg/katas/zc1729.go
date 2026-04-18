package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1729",
		Title:    "Error on `ip route flush all` / `ip route del default` — script loses network connectivity",
		Severity: SeverityError,
		Description: "`ip route flush all` (or `flush table main`) wipes every routing entry, " +
			"including the default gateway. `ip route del default` removes only the default " +
			"route — same outcome. The remote SSH session that just ran the command can " +
			"no longer talk to the host, and any subsequent step that needs the network " +
			"hangs until manual console intervention. Scope the flush (`flush dev <iface>`, " +
			"`flush scope link`) or use `ip route replace default via <gw>` so the new " +
			"route is in place before the old one disappears.",
		Check: checkZC1729,
	})
}

func checkZC1729(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ip" {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, arg := range cmd.Arguments {
		v := arg.String()
		// Skip leading short flags like `-4`, `-6`, `-s`, `-d`.
		if len(args) == 0 && strings.HasPrefix(v, "-") {
			continue
		}
		args = append(args, v)
	}

	if len(args) < 3 || args[0] != "route" {
		return nil
	}

	switch args[1] {
	case "flush":
		if args[2] == "all" || args[2] == "table" {
			return zc1729Hit(cmd, "ip route flush "+args[2])
		}
	case "del", "delete":
		if args[2] == "default" {
			return zc1729Hit(cmd, "ip route "+args[1]+" default")
		}
	}
	return nil
}

func zc1729Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1729",
		Message: "`" + what + "` removes the default gateway — the SSH session that " +
			"just ran it loses connectivity. Scope the flush (`flush dev <iface>`) or " +
			"use `ip route replace default via <gw>` so the new route is in place " +
			"first.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
