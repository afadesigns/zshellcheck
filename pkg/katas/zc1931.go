// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1931",
		Title:    "Warn on `ip netns delete $NS` / `netns del` — drops the whole network namespace",
		Severity: SeverityWarning,
		Description: "`ip netns delete NAME` / `ip netns del NAME` unmounts the namespace and " +
			"tears down every interface, veth pair, VXLAN, and WireGuard peer living inside. " +
			"Processes still attached lose their network abruptly — container health checks " +
			"fail, BGP sessions drop, and any other process using `ip netns exec NAME …` " +
			"errors out with \"No such file or directory\". Stop the workloads first " +
			"(`systemctl stop`, `pkill -SIGTERM -n $NS`), confirm `ip -n $NS link` is empty, " +
			"then `delete` deliberately — or leave the namespace alone if it is managed by " +
			"Docker/containerd/systemd-nspawn.",
		Check: checkZC1931,
	})
}

func checkZC1931(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ip" {
		return nil
	}
	if len(cmd.Arguments) < 3 {
		return nil
	}
	if cmd.Arguments[0].String() != "netns" {
		return nil
	}
	sub := cmd.Arguments[1].String()
	if sub != "delete" && sub != "del" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1931",
		Message: "`ip netns " + sub + "` tears down every interface, veth, tunnel, and " +
			"WireGuard peer inside the namespace. Stop the workloads first and verify " +
			"`ip -n $NS link` is empty before deleting.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
