package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC2002",
		Title:    "Error on `crictl rmi -a` / `crictl rm -af` — wipes every image/container on the Kubernetes node",
		Severity: SeverityError,
		Description: "`crictl` talks directly to the node's CRI runtime (containerd, CRI-O), " +
			"below the kubelet and the cluster API. `crictl rmi -a` removes every " +
			"cached image including the ones currently backing running pods — the " +
			"kubelet must immediately re-pull from the registry, and image-pull rate " +
			"limits or network blips turn the node Unready. `crictl rm -af` force-" +
			"removes every container on the node, killing pods without running " +
			"PreStop hooks or honoring PodDisruptionBudget. Route maintenance through " +
			"`kubectl drain $NODE` + `kubectl delete pod --grace-period=30`; use " +
			"`crictl` at most on a cordoned, drained node with a documented recovery " +
			"plan.",
		Check: checkZC2002,
	})
}

func checkZC2002(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "crictl" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "rmi" && sub != "rm" {
		return nil
	}
	hasAll := false
	hasForce := false
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "-a" || v == "--all" || v == "-af" || v == "-fa" {
			hasAll = true
		}
		if v == "-f" || v == "--force" || v == "-af" || v == "-fa" {
			hasForce = true
		}
	}
	if sub == "rmi" && hasAll {
		return zc2002Hit(cmd, "crictl rmi -a")
	}
	if sub == "rm" && hasAll && hasForce {
		return zc2002Hit(cmd, "crictl rm -af")
	}
	return nil
}

func zc2002Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC2002",
		Message: "`" + form + "` talks to the node CRI directly, under the kubelet — " +
			"images/containers backing running pods disappear, kubelet must " +
			"re-pull or re-run. Route through `kubectl drain`/`delete pod`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}
