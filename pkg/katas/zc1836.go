package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1836",
		Title:    "Error on `helm uninstall --no-hooks` — skips pre-delete cleanup, orphaned state",
		Severity: SeverityError,
		Description: "`helm uninstall RELEASE --no-hooks` (also spelled `helm delete --no-hooks` on " +
			"Helm v2 / `helm3 --no-hooks` interchangeably) tears down every chart-rendered " +
			"resource but silently skips the release's `pre-delete` and `post-delete` " +
			"Jobs / ConfigMap hooks. Those hooks are where production charts flush " +
			"write-ahead logs, deregister service-discovery entries, back up PVC content " +
			"before the PVC goes away, and release external locks — skipping them on a " +
			"live release is one of the classic ways to leave the cluster in a partially " +
			"deleted state with no way to replay the cleanup. Drop `--no-hooks` and let " +
			"the chart run as designed; if a hook is genuinely wedged, disable it at the " +
			"chart level with `helm.sh/hook-delete-policy: before-hook-creation,hook-succeeded`.",
		Check: checkZC1836,
	})
}

func checkZC1836(node ast.Node) []Violation {
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
	args := cmd.Arguments
	if len(args) < 2 {
		return nil
	}
	sub := args[0].String()
	if sub != "uninstall" && sub != "delete" {
		return nil
	}
	for _, arg := range args[1:] {
		if arg.String() == "--no-hooks" {
			return []Violation{{
				KataID: "ZC1836",
				Message: "`helm " + sub + " --no-hooks` skips pre/post-delete " +
					"cleanup hooks — orphaned locks, DNS, missed PVC backups. " +
					"Drop the flag; fix stuck hooks via " +
					"`helm.sh/hook-delete-policy`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
