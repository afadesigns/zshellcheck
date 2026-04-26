// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1547",
		Title:    "Warn on `kubectl apply --prune --all` — deletes resources missing from manifest",
		Severity: SeverityWarning,
		Description: "`kubectl apply --prune --all` (or `--prune -l <selector>`) deletes every " +
			"cluster resource whose label matches but which is not in the manifest you just " +
			"applied. In a partial-repo deploy or a manifest typo, that can delete production " +
			"Deployments, Services, or Secrets another team owns. Pair `--prune` with a " +
			"narrow `-l` selector unique to your stack, or use a GitOps controller (Argo CD, " +
			"Flux) that scopes prune to its own Application.",
		Check: checkZC1547,
	})
}

func checkZC1547(node ast.Node) []Violation {
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

	var sawApply, hasPrune, hasAll bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "apply" {
			sawApply = true
		}
		if !sawApply {
			continue
		}
		if v == "--prune" {
			hasPrune = true
		}
		if v == "--all" || v == "-A" || v == "--all-namespaces" {
			hasAll = true
		}
	}
	if sawApply && hasPrune && hasAll {
		return []Violation{{
			KataID: "ZC1547",
			Message: "`kubectl apply --prune --all` deletes every matching resource not in the " +
				"manifest — manifest typo wipes other teams' resources. Scope with a " +
				"narrow `-l <selector>`.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}
