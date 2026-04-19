package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1878",
		Title:    "Warn on `kubectl apply --force-conflicts` — steals ownership of fields managed by other controllers",
		Severity: SeverityWarning,
		Description: "Server-side apply tracks every field of a resource by the applier that " +
			"last set it (`metadata.managedFields`). When two appliers disagree, the " +
			"default behaviour is to abort with `conflict` so you can reconcile " +
			"deliberately. `kubectl apply --server-side --force-conflicts` overrides " +
			"that: the current caller snatches ownership of every conflicting field — " +
			"including fields set by operators, HPA, cert-manager, and webhook-injected " +
			"sidecars — and those controllers will silently lose their reconcile " +
			"pressure until their next write. Resolve the conflict instead: either " +
			"drop the disputed fields from your manifest so the other owner can keep " +
			"them, or coordinate a hand-off by first removing the managed-field entry " +
			"(`kubectl apply --field-manager=... --subresource=...`).",
		Check: checkZC1878,
	})
}

func checkZC1878(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "kubectl" {
		return nil
	}
	args := cmd.Arguments
	if len(args) == 0 || args[0].String() != "apply" {
		return nil
	}
	for _, arg := range args[1:] {
		if arg.String() == "--force-conflicts" {
			return []Violation{{
				KataID: "ZC1878",
				Message: "`kubectl apply --force-conflicts` grabs ownership of every " +
					"conflicting field from other controllers (HPA, cert-manager, " +
					"sidecar injectors). Resolve the conflict instead — drop the " +
					"disputed fields or hand off via managed-field edit.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
