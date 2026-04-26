// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1762",
		Title:    "Error on `kubeadm join --discovery-token-unsafe-skip-ca-verification` — cluster CA not checked",
		Severity: SeverityError,
		Description: "`kubeadm join` verifies the control-plane API server's CA before accepting " +
			"the kubelet bootstrap token. `--discovery-token-unsafe-skip-ca-verification` " +
			"skips that check, so a network-position attacker can impersonate the API " +
			"server, harvest the bootstrap token, and seed malicious workloads onto the " +
			"joining node. Always pin the CA with `--discovery-token-ca-cert-hash sha256:" +
			"<digest>` (emitted by `kubeadm token create --print-join-command`) or supply " +
			"a kubeconfig discovery file that has the CA baked in.",
		Check: checkZC1762,
	})
}

func checkZC1762(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "kubeadm" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "join" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		if arg.String() == "--discovery-token-unsafe-skip-ca-verification" {
			return []Violation{{
				KataID: "ZC1762",
				Message: "`kubeadm join --discovery-token-unsafe-skip-ca-verification` " +
					"skips CA verification of the control-plane — MITM steals the " +
					"bootstrap token. Pin the CA with `--discovery-token-ca-cert-hash " +
					"sha256:<digest>`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
