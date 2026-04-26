// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1880",
		Title:    "Warn on `kubectl annotate|label --overwrite` — silently rewrites controller signals",
		Severity: SeverityWarning,
		Description: "Kubernetes annotations and labels are not plain metadata — they are the " +
			"protocol by which cert-manager, external-dns, ingress-nginx, the " +
			"HorizontalPodAutoscaler, and most Helm-managed controllers decide what to " +
			"do with a resource. `kubectl annotate --overwrite` and `kubectl label " +
			"--overwrite` suppress the conflict check and replace whatever value was " +
			"there, so the script silently rewrites `kubectl.kubernetes.io/last-applied-" +
			"configuration`, `cert-manager.io/cluster-issuer`, or " +
			"`prometheus.io/scrape`, triggering reissue / reconfiguration or breaking " +
			"the next apply. Inspect the existing annotation with `kubectl get -o " +
			"jsonpath='{.metadata.annotations}'` first, and drop `--overwrite` so a " +
			"conflict surfaces as an error.",
		Check: checkZC1880,
	})
}

func checkZC1880(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "kubectl" {
		return nil
	}
	args := cmd.Arguments
	if len(args) == 0 {
		return nil
	}
	sub := args[0].String()
	if sub != "annotate" && sub != "label" {
		return nil
	}
	for _, arg := range args[1:] {
		if arg.String() == "--overwrite" {
			return []Violation{{
				KataID: "ZC1880",
				Message: "`kubectl " + sub + " --overwrite` silently replaces an " +
					"existing controller signal — cert-manager, external-dns, " +
					"HPA watchers reconcile on the new value. Inspect first; " +
					"drop `--overwrite` so conflicts error.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
