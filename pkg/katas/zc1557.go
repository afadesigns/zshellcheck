// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1557",
		Title:    "Error on `kubeadm reset -f` / `--force` — wipes Kubernetes control-plane state",
		Severity: SeverityError,
		Description: "`kubeadm reset` stops kubelet, tears down static-pod manifests, clears " +
			"`/etc/kubernetes`, and (with `-f`) skips the confirmation that protects a mistyped " +
			"target. On a control-plane node it also breaks every tenant that relied on that " +
			"etcd quorum. Drain first, remove the node from the cluster, then run reset " +
			"interactively to confirm.",
		Check: checkZC1557,
	})
}

func checkZC1557(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kubeadm" {
		return nil
	}

	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "reset" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "-f" || v == "--force" {
			return []Violation{{
				KataID: "ZC1557",
				Message: "`kubeadm reset -f` skips the confirmation and wipes " +
					"/etc/kubernetes / kubelet state. Drain and remove the node first.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
