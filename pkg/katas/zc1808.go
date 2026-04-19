package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1808",
		Title:    "Warn on `kubectl replace --force` — deletes + recreates resource, drops running pods",
		Severity: SeverityWarning,
		Description: "`kubectl replace --force -f FILE` is `delete` followed by `create`: the " +
			"existing resource (and every dependent pod / replicaset / endpoint) is removed " +
			"before the new manifest is applied. In-flight requests drop, PodDisruptionBudget " +
			"is ignored, and controllers that watch the object see it disappear and reappear. " +
			"Prefer `kubectl apply -f FILE` — same manifest, server-side merge that preserves " +
			"running pods — and reach for `replace --force` only when the resource schema has " +
			"changed in a way `apply` cannot patch, with traffic drained beforehand.",
		Check: checkZC1808,
	})
}

func checkZC1808(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "kubectl" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "replace" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--force" || v == "-f=--force" {
			return []Violation{{
				KataID: "ZC1808",
				Message: "`kubectl replace --force` is delete + create — pods die, " +
					"PDBs are ignored, in-flight requests drop. Prefer `kubectl " +
					"apply -f FILE` and reserve `replace --force` for schema changes " +
					"`apply` cannot patch, after draining traffic.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
