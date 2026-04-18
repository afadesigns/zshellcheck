package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1744",
		Title:    "Warn on `kubectl port-forward --address 0.0.0.0` — cluster port exposed to every interface",
		Severity: SeverityWarning,
		Description: "`kubectl port-forward` defaults to binding the local end of the tunnel on " +
			"`127.0.0.1`. `--address 0.0.0.0` (or a specific non-loopback IP) exposes the " +
			"target pod's port to every interface on the developer's workstation or the " +
			"bastion host running the command. Anyone on the LAN / VPN can reach internal " +
			"cluster services that never meant to be externally reachable. Drop the flag " +
			"(loopback default), or pick a specific interface that is already scoped to a " +
			"trusted network.",
		Check: checkZC1744,
	})
}

func checkZC1744(node ast.Node) []Violation {
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
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "port-forward" {
		return nil
	}

	prevAddress := false
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if prevAddress {
			if v == "0.0.0.0" || v == "::" {
				return zc1744Hit(cmd, "--address "+v)
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
				return zc1744Hit(cmd, v)
			}
		}
	}
	return nil
}

func zc1744Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1744",
		Message: "`kubectl port-forward " + what + "` binds the local end of the tunnel " +
			"on every interface — anyone on the LAN / VPN can reach the pod. Drop " +
			"`--address` (loopback default) or pick a trusted-network interface IP.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
