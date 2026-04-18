package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1676",
		Title:    "Warn on `helm rollback --force` — recreates in-flight resources, corrupts rolling updates",
		Severity: SeverityWarning,
		Description: "`helm rollback RELEASE N --force` asks Helm to delete and recreate any " +
			"resource that it cannot patch cleanly. If a deployment is mid-rollout, the " +
			"`--force` flag takes out both the old and new ReplicaSets, kicks the pods, " +
			"and forces a cold start — losing in-flight requests and any `PodDisruptionBudget` " +
			"protections. Worse, rolling back to revision N brings back whatever CVEs or " +
			"config regressions the later revisions had already fixed. Pin the target " +
			"revision explicitly, omit `--force`, and gate the rollback behind a change-" +
			"review ticket rather than a shell one-liner.",
		Check: checkZC1676,
	})
}

func checkZC1676(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "helm" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "rollback" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		if arg.String() == "--force" {
			return []Violation{{
				KataID: "ZC1676",
				Message: "`helm rollback --force` deletes and recreates unpatched resources — " +
					"loses in-flight traffic and bypasses PodDisruptionBudget. Drop `--force` " +
					"and gate the rollback via change review.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
