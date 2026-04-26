// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1933",
		Title:    "Error on `ipvsadm -C` / `--clear` — wipes every IPVS virtual service, drops load balancer",
		Severity: SeverityError,
		Description: "`ipvsadm -C` (and the long form `--clear`) removes every virtual service, " +
			"real server, and connection entry from the in-kernel IPVS table. Traffic that was " +
			"being load-balanced to a backend farm now falls through to the host's local " +
			"listen sockets (or drops), active keepalived/`ldirectord` states invert, and " +
			"clients see 5xx until an operator replays the config. Save the current table first " +
			"(`ipvsadm --save -n > /run/ipvs.bak`), drain specific services with `ipvsadm -D`, " +
			"and keep `--clear` in break-glass-only runbooks.",
		Check: checkZC1933,
	})
}

var zc1933ClearFlags = map[string]bool{
	"-C":      true,
	"--clear": true,
}

func checkZC1933(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ipvsadm" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-C" || v == "--clear" {
			return zc1933Hit(cmd, "ipvsadm "+v)
		}
	}
	return nil
}

func zc1933Hit(cmd *ast.SimpleCommand, form string) []Violation {
	line, col := FlagArgPosition(cmd, zc1933ClearFlags)
	return []Violation{{
		KataID: "ZC1933",
		Message: "`" + form + "` wipes every IPVS virtual service and real-server binding — " +
			"load balancing stops, clients see 5xx. Save via `ipvsadm --save`, drain " +
			"specific services with `-D`, reserve `--clear` for break-glass.",
		Line:   line,
		Column: col,
		Level:  SeverityError,
	}}
}
