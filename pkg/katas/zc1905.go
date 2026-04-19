package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1905",
		Title:    "Warn on `ssh -g -L …` — local forward bound on all interfaces, not just loopback",
		Severity: SeverityWarning,
		Description: "`ssh -g` flips the default for `-L` (local forward) and `-D` (dynamic SOCKS) " +
			"from `127.0.0.1:port` to `0.0.0.0:port`. Any host on the same LAN/VPN/WiFi " +
			"segment can then use the tunnel without authenticating to the SSH session. " +
			"Drop `-g`, pin the bind explicitly with `-L bind_address:port:target:port`, or " +
			"use a firewall rule — never leave a forwarded port open to the network segment.",
		Check: checkZC1905,
	})
}

func checkZC1905(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ssh" {
		return nil
	}

	hasG := false
	hasForward := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch {
		case v == "-g":
			hasG = true
		case strings.HasPrefix(v, "-g") && len(v) > 2 && !strings.HasPrefix(v, "-go"):
			// clustered short flags, e.g. `-gNL`
			hasG = true
		case v == "-L", v == "-D":
			hasForward = true
		case strings.HasPrefix(v, "-L") && len(v) > 2:
			hasForward = true
		case strings.HasPrefix(v, "-D") && len(v) > 2:
			hasForward = true
		}
		if strings.HasPrefix(v, "-") && len(v) >= 2 && v[1] != '-' {
			for i := 1; i < len(v); i++ {
				if v[i] == 'g' {
					hasG = true
				}
				if v[i] == 'L' || v[i] == 'D' {
					hasForward = true
				}
			}
		}
	}
	if !hasG || !hasForward {
		return nil
	}

	return []Violation{{
		KataID: "ZC1905",
		Message: "`ssh -g` with `-L`/`-D` binds the forward on `0.0.0.0` — anyone on the " +
			"same LAN segment can ride the tunnel. Drop `-g` or pin `bind_address:port` " +
			"in the forward spec.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
