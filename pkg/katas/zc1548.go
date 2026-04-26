// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1548",
		Title:    "Warn on `helm install/upgrade --disable-openapi-validation` — skips schema check",
		Severity: SeverityWarning,
		Description: "`--disable-openapi-validation` tells Helm to skip the OpenAPI schema check " +
			"the API server would apply. Malformed CRD instances or Deployments with " +
			"invalid spec fields then silently land in etcd, only failing when the " +
			"controller tries to reconcile — usually 3am, usually in prod. Keep the " +
			"validation on; fix the schema deviation instead.",
		Check: checkZC1548,
	})
}

func checkZC1548(node ast.Node) []Violation {
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

	var sawVerb bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "install" || v == "upgrade" || v == "template" {
			sawVerb = true
			continue
		}
		if !sawVerb {
			continue
		}
		if v == "--disable-openapi-validation" {
			return []Violation{{
				KataID: "ZC1548",
				Message: "`helm --disable-openapi-validation` hides bad manifests until the " +
					"controller crashes. Fix the schema deviation.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
