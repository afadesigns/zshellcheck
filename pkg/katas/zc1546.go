// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1546",
		Title:    "Warn on `kubectl delete --force --grace-period=0` — skips PreStop, corrupts state",
		Severity: SeverityWarning,
		Description: "`kubectl delete --force --grace-period=0` tells the API server to remove " +
			"the resource from etcd without waiting for the kubelet to run PreStop hooks or " +
			"drain the pod. For a StatefulSet pod this routinely corrupts the backing PV " +
			"(database mid-flush, file lock left held) and the replacement pod refuses to " +
			"start. Use standard delete and let the graceful shutdown run; only reach for " +
			"`--force` when the node itself is gone.",
		Check: checkZC1546,
	})
}

func checkZC1546(node ast.Node) []Violation {
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

	var sawDelete, hasForce, hasGrace0 bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "delete" {
			sawDelete = true
		}
		if !sawDelete {
			continue
		}
		if v == "--force" {
			hasForce = true
		}
		if v == "--grace-period=0" {
			hasGrace0 = true
		}
	}
	if sawDelete && hasForce && hasGrace0 {
		return []Violation{{
			KataID: "ZC1546",
			Message: "`kubectl delete --force --grace-period=0` skips PreStop hooks and kubelet " +
				"drain — corrupts StatefulSet state. Use standard delete.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}
