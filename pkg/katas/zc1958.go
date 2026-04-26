// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1958",
		Title:    "Warn on `helm upgrade --force` — delete-and-recreate resources, drops running pods",
		Severity: SeverityWarning,
		Description: "`helm upgrade RELEASE CHART --force` flips the upgrade strategy from " +
			"three-way-merge to `delete + create` for every resource Helm owns. Deployments " +
			"become new objects, Services lose their `clusterIP` for a beat, and any " +
			"`PodDisruptionBudget` is bypassed because the resource is deleted, not rolled " +
			"out. Use plain `helm upgrade` (three-way merge) or `--atomic` / `--wait` for a " +
			"supervised roll. Reserve `--force` for recovery after a failed upgrade with a " +
			"stuck resource, not routine deploys.",
		Check: checkZC1958,
	})
}

func checkZC1958(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "helm" && ident.Value != "helm3" {
		return nil
	}
	if len(cmd.Arguments) < 2 || cmd.Arguments[0].String() != "upgrade" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		if arg.String() == "--force" {
			return zc1958Hit(cmd)
		}
	}
	return nil
}

func zc1958Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1958",
		Message: "`helm upgrade --force` is delete+create — pods die, PodDisruptionBudget " +
			"is bypassed, Services reset their `clusterIP`. Use plain `helm upgrade` " +
			"(three-way merge) or `--atomic`/`--wait` for a supervised roll.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
