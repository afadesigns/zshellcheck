// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1750",
		Title:    "Error on `kubectl proxy --address 0.0.0.0` — cluster API proxy on every interface",
		Severity: SeverityError,
		Description: "`kubectl proxy` tunnels Kubernetes API requests authenticated with the " +
			"local kubeconfig's credentials. Defaults bind to `127.0.0.1` and accept only " +
			"`localhost` hosts. `--address 0.0.0.0` (or a specific non-loopback IP) exposes " +
			"that tunnel to every interface on the workstation / bastion, so anyone on the " +
			"LAN or VPN gets the cluster admin the kubeconfig holds. Same risk applies to " +
			"`--accept-hosts '.*'`. Keep the loopback default and scope with SSH port " +
			"forwarding, or restrict `--address` to an interface behind a tight firewall.",
		Check: checkZC1750,
	})
}

func checkZC1750(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kubectl" && ident.Value != "oc" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "proxy" {
		return nil
	}

	prevAddress := false
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if prevAddress {
			if v == "0.0.0.0" || v == "::" {
				return zc1750Hit(cmd, "--address "+v)
			}
			prevAddress = false
			continue
		}
		switch {
		case v == "--address":
			prevAddress = true
		case strings.HasPrefix(v, "--address="):
			val := strings.TrimPrefix(v, "--address=")
			if val == "0.0.0.0" || val == "::" {
				return zc1750Hit(cmd, v)
			}
		}
	}
	return nil
}

func zc1750Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1750",
		Message: "`kubectl proxy " + what + "` exposes the cluster-admin API tunnel to every " +
			"reachable interface. Keep the loopback default and tunnel over SSH, or " +
			"restrict `--address` to a firewalled interface.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
