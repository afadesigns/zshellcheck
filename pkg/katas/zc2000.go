package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC2000",
		Title:    "Error on `kubectl taint nodes $NODE key=value:NoExecute` — evicts every non-tolerating pod off the node",
		Severity: SeverityError,
		Description: "A `NoExecute` taint kicks every existing pod off the node unless the pod " +
			"spec explicitly tolerates it. Draining one node during a rolling upgrade " +
			"is one thing; a script that types the taint wrong (mis-keying the " +
			"toleration value, applying to `--all` nodes, or iterating a node list " +
			"without a pause) can empty a whole cluster in seconds and trigger " +
			"cascade reschedules that overwhelm the scheduler. Prefer `kubectl drain " +
			"$NODE` (which respects PodDisruptionBudget and runs PreStop hooks) or a " +
			"`NoSchedule` taint for gentle drain; reserve `NoExecute` for genuine " +
			"incident response with a runbook and a safety countdown.",
		Check: checkZC2000,
	})
}

func checkZC2000(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kubectl" {
		return nil
	}
	if len(cmd.Arguments) < 2 || cmd.Arguments[0].String() != "taint" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if strings.Contains(v, ":NoExecute") {
			return []Violation{{
				KataID: "ZC2000",
				Message: "`kubectl taint nodes … :NoExecute` evicts every non-tolerating " +
					"pod immediately — a typo on `--all` nodes empties the cluster. " +
					"Prefer `kubectl drain $NODE` or a `:NoSchedule` taint for " +
					"gentle drain.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
