// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1471",
		Title:    "Error on `kubectl/helm --insecure-skip-tls-verify` (cluster MITM)",
		Severity: SeverityError,
		Description: "`--insecure-skip-tls-verify` tells kubectl / helm to accept any certificate " +
			"from the API server. Against a production cluster, this hands every secret and " +
			"admission payload to a MITM. Fix the trust chain: point `--certificate-authority` " +
			"at the right CA bundle, or restore `KUBECONFIG` with the cluster's embedded CA.",
		Check: checkZC1471,
	})
}

func checkZC1471(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kubectl" && ident.Value != "helm" && ident.Value != "oc" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--insecure-skip-tls-verify" ||
			v == "--insecure-skip-tls-verify=true" ||
			v == "--kube-insecure-skip-tls-verify" {
			return []Violation{{
				KataID: "ZC1471",
				Message: "`--insecure-skip-tls-verify` turns off API-server certificate " +
					"verification — MITM steals every secret. Fix the CA bundle instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
