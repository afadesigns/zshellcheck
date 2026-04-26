// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1711",
		Title:    "Error on `etcdctl del --prefix \"\"` / `--from-key \"\"` — wipes the entire keyspace",
		Severity: SeverityError,
		Description: "`etcdctl del --prefix KEY` deletes every key under KEY's range. With KEY " +
			"empty (`\"\"` or `\"\\0\"`) the range is `[\"\", \"\\xFF\")` — the whole etcd " +
			"cluster, including kube-apiserver state if etcd is the Kubernetes datastore. " +
			"`--from-key \"\"` has the same effect for the lower-bound form. Restrict the " +
			"prefix to the namespace you actually own (`/app/staging/`), or wrap the call " +
			"with an explicit `etcdctl get --prefix --keys-only` review step.",
		Check: checkZC1711,
	})
}

func checkZC1711(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "etcdctl" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "del" {
		return nil
	}

	for i, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v != "--prefix" && v != "--from-key" {
			continue
		}
		idx := i + 2
		if idx >= len(cmd.Arguments) {
			continue
		}
		next := cmd.Arguments[idx].String()
		if next == `""` || next == "''" {
			return []Violation{{
				KataID: "ZC1711",
				Message: "`etcdctl del " + v + " \"\"` deletes the entire etcd keyspace " +
					"(including kube-apiserver state) — scope to a specific namespace " +
					"prefix and review with `get --prefix --keys-only` first.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
